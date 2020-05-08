// Copyright 2019 Netflix, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package experimentrouter

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/Netflix/p2plab/nodes"
	"github.com/Netflix/p2plab/reports"
	"github.com/uber/jaeger-client-go"
	bolt "go.etcd.io/bbolt"

	"github.com/Netflix/p2plab"
	"github.com/Netflix/p2plab/daemon"
	"github.com/Netflix/p2plab/labd/controlapi"
	"github.com/Netflix/p2plab/labd/routers/helpers"
	"github.com/Netflix/p2plab/metadata"
	"github.com/Netflix/p2plab/peer"
	"github.com/Netflix/p2plab/pkg/httputil"
	"github.com/Netflix/p2plab/pkg/stringutil"
	"github.com/Netflix/p2plab/query"
	"github.com/Netflix/p2plab/scenarios"
	"github.com/Netflix/p2plab/transformers"
	"github.com/containerd/containerd/errdefs"
	"github.com/pkg/errors"
	"github.com/rs/xid"
	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"
)

type router struct {
	db       metadata.DB
	provider p2plab.NodeProvider
	client   *httputil.Client
	ts       *transformers.Transformers
	seeder   *peer.Peer
	builder  p2plab.Builder
	rhelper  *helpers.Helper
}

// New returns a new experiment router initialized with the router helpers
func New(db metadata.DB, provider p2plab.NodeProvider, client *httputil.Client, ts *transformers.Transformers, seeder *peer.Peer, builder p2plab.Builder) daemon.Router {
	return &router{
		db,
		provider,
		client,
		ts,
		seeder,
		builder,
		helpers.New(db, provider, client),
	}
}

func (s *router) Routes() []daemon.Route {
	return []daemon.Route{
		// GET
		daemon.NewGetRoute("/experiments/json", s.getExperiments),
		daemon.NewGetRoute("/experiments/{id}/json", s.getExperimentByName),
		// POST
		daemon.NewPostRoute("/experiments/create", s.postExperimentsCreate),
		// PUT
		daemon.NewPutRoute("/experiments/label", s.putExperimentsLabel),
		// DELETE
		daemon.NewDeleteRoute("/experiments/delete", s.deleteExperiments),
	}
}

func (s *router) getExperiments(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	experiments, err := s.db.ListExperiments(ctx)
	if err != nil {
		return err
	}

	return daemon.WriteJSON(w, &experiments)
}

func (s *router) getExperimentByName(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	id := vars["name"]
	experiment, err := s.db.GetExperiment(ctx, id)
	if err != nil {
		return err
	}

	return daemon.WriteJSON(w, &experiment)
}

func (s *router) postExperimentsCreate(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	noReset := false
	if r.FormValue("no-reset") != "" {
		var err error
		noReset, err = strconv.ParseBool(r.FormValue("no-reset"))
		if err != nil {
			return err
		}
	}

	defer r.Body.Close()
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	var edef metadata.ExperimentDefinition
	if err := edef.FromJSON(data); err != nil {
		return err
	}
	exp, err := s.db.CreateExperiment(ctx, metadata.Experiment{
		ID:         xid.New().String(),
		Definition: edef,
		Status:     metadata.ExperimentRunning,
	})
	if err != nil {
		return err
	}
	errg, ctx := errgroup.WithContext(ctx)
	for i, trial := range exp.Definition.TrialDefinition {
		// TODO(bonedaddy): if we dont do this, then we run into some issues
		// with port numbers being re-used. For example when we start the goroutines
		// for each of the benchmarks they use the same ports for all peers
		// and as such we run into port issues
		if i > 0 {
			break
		}
		trial := trial
		name := fmt.Sprintf("%s-%v", xid.New().String(), i)
		errg.Go(func() error {
			info := zerolog.Ctx(ctx).Info()
			cluster, err := s.rhelper.CreateCluster(ctx, trial.Cluster, name, w)
			if err != nil {
				return err
			}
			info.Msg("creating scenario")
			scenID := xid.New().String()
			scenario := metadata.Scenario{
				ID:         scenID,
				Definition: trial.Scenario,
				Labels: []string{
					name,
				},
			}
			scenario, err = s.db.CreateScenario(ctx, scenario)
			if err != nil {
				return err
			}
			defer func() error {
				info.Msg("tearing down cluster")
				cluster, err := s.db.GetCluster(ctx, name)
				if err != nil {
					return errors.Wrap(err, "failed to get cluster'")
				}
				ns, err := s.db.ListNodes(ctx, cluster.ID)
				if err != nil {
					return errors.Wrap(err, "failed to list nodes")
				}
				ng := &p2plab.NodeGroup{
					ID:    cluster.ID,
					Nodes: ns,
				}
				if err := s.provider.DestroyNodeGroup(ctx, ng); err != nil {
					return errors.Wrap(err, "failed to destroy node group")
				}
				info.Msg("tore down cluster")
				return nil
			}()
			info.Msg("creating nodes")
			mns, err := s.db.ListNodes(ctx, cluster.ID)
			if err != nil {
				return err
			}
			var (
				ns   []p2plab.Node
				lset = query.NewLabeledSet()
			)
			for _, n := range mns {
				node := controlapi.NewNode(s.client, n)
				lset.Add(node)
				ns = append(ns, node)
			}
			if !noReset {
				info.Msg("updating nodes")
				if err := nodes.Update(ctx, s.builder, ns); err != nil {
					return errors.Wrap(err, "failed to update cluster")
				}
				info.Msg("connecting nodes")
				if err := nodes.Connect(ctx, ns); err != nil {
					return errors.Wrap(err, "failed to connect cluster")
				}
			}
			info.Msg("generating scenario plan")
			plan, queries, err := scenarios.Plan(ctx, trial.Scenario, s.ts, s.seeder, lset)
			if err != nil {
				return err
			}
			bid := xid.New().String()
			benchmark := metadata.Benchmark{
				ID:       bid,
				Status:   metadata.BenchmarkRunning,
				Cluster:  cluster,
				Scenario: scenario,
				Plan:     plan,
				Labels: []string{
					cluster.ID,
					scenario.ID,
					bid,
				},
			}
			info.Msg("creating benchmark")
			if benchmark, err = s.db.CreateBenchmark(ctx, benchmark); err != nil {
				return err
			}
			var seederAddrs []string
			for _, addr := range s.seeder.Host().Addrs() {
				seederAddrs = append(seederAddrs, fmt.Sprintf("%s/p2p/%s", addr, s.seeder.Host().ID()))
			}
			info.Msg("running scenario")
			execution, err := scenarios.Run(ctx, lset, plan, seederAddrs)
			if err != nil {
				return errors.Wrap(err, "failed to run scenario plan")
			}
			report := metadata.Report{
				Summary: metadata.ReportSummary{
					TotalTime: execution.End.Sub(execution.Start),
				},
				Nodes:   execution.Report,
				Queries: queries,
			}
			info.Msg("aggregating results")
			report.Aggregates = reports.ComputeAggregates(report.Nodes)
			jaegerUI := os.Getenv("JAEGER_UI")
			if jaegerUI != "" {
				sc, ok := execution.Span.Context().(jaeger.SpanContext)
				if ok {
					report.Summary.Trace = fmt.Sprintf("%s/trace/%s", jaegerUI, sc.TraceID())
				}
			}
			info.Msg("Updating benchmark metadata")
			err = s.db.Update(ctx, func(tx *bolt.Tx) error {
				tctx := metadata.WithTransactionContext(ctx, tx)

				err := s.db.CreateReport(tctx, benchmark.ID, report)
				if err != nil {
					return errors.Wrap(err, "failed to create report")
				}

				benchmark.Status = metadata.BenchmarkDone
				_, err = s.db.UpdateBenchmark(tctx, benchmark)
				if err != nil {
					return errors.Wrap(err, "failed to update benchmark")
				}

				return nil
			})
			info.Msg("updating reports")
			exp.Reports = append(exp.Reports, report)
			return err
		})
	}
	err = errg.Wait()
	daemon.WriteJSON(w, &exp)
	return err
}

func (s *router) putExperimentsLabel(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	ids := strings.Split(r.FormValue("ids"), ",")
	addLabels := stringutil.Coalesce(strings.Split(r.FormValue("adds"), ","))
	removeLabels := stringutil.Coalesce(strings.Split(r.FormValue("removes"), ","))

	var experiments []metadata.Experiment
	if len(addLabels) > 0 || len(removeLabels) > 0 {
		var err error
		experiments, err = s.db.LabelExperiments(ctx, ids, addLabels, removeLabels)
		if err != nil {
			return err
		}
	}

	return daemon.WriteJSON(w, &experiments)
}

func (s *router) deleteExperiments(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	ids := strings.Split(r.FormValue("ids"), ",")

	for _, id := range ids {
		logger := zerolog.Ctx(ctx).With().Str("experiment", id).Logger()

		experiment, err := s.db.GetExperiment(ctx, id)
		if err != nil {
			return err
		}

		switch experiment.Status {
		case metadata.ExperimentDone, metadata.ExperimentError:
		default:
			return errors.Wrapf(errdefs.ErrFailedPrecondition, "experiment status %q", experiment.Status)
		}

		logger.Info().Msg("Deleting experiment")
		err = s.db.DeleteExperiment(ctx, id)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *router) matchExperiments(ctx context.Context, q string) ([]metadata.Experiment, error) {
	es, err := s.db.ListExperiments(ctx)
	if err != nil {
		return nil, err
	}

	var ls []p2plab.Labeled
	for _, e := range es {
		ls = append(ls, query.NewLabeled(e.ID, e.Labels))
	}

	mset, err := query.Execute(ctx, ls, q)
	if err != nil {
		return nil, err
	}

	var matchedExperiments []metadata.Experiment
	for _, e := range es {
		if mset.Contains(e.ID) {
			matchedExperiments = append(matchedExperiments, e)
		}
	}

	return matchedExperiments, nil
}
