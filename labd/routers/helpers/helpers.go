package helpers

import (
	"context"
	"net/http"

	"github.com/Netflix/p2plab"
	"github.com/Netflix/p2plab/labd/controlapi"
	"github.com/Netflix/p2plab/metadata"
	"github.com/Netflix/p2plab/nodes"
	"github.com/Netflix/p2plab/pkg/httputil"
	"github.com/rs/zerolog"
	bolt "go.etcd.io/bbolt"
)

// TODO(bonedaddy): not sure if this is the best way to go about sharing code between routers

// Helper abstracts commonly used functions to be shared by any router
type Helper struct {
	db       metadata.DB
	provider p2plab.NodeProvider
	client   *httputil.Client
}

// New instantiates our helper type
func New(db metadata.DB, provider p2plab.NodeProvider, client *httputil.Client) *Helper {
	return &Helper{db, provider, client}
}

// CreateCluster enables creating the nodes in a cluster, waiting for them to be healthy before returning
func (h *Helper) CreateCluster(ctx context.Context, cdef metadata.ClusterDefinition, name string, w http.ResponseWriter) (metadata.Cluster, error) {
	var (
		cluster = metadata.Cluster{
			ID:         name,
			Status:     metadata.ClusterCreating,
			Definition: cdef,
			Labels: append([]string{
				name,
			}, cdef.GenerateLabels()...),
		}
		err error
	)

	cluster, err = h.db.CreateCluster(ctx, cluster)
	if err != nil {
		return cluster, err
	}

	zerolog.Ctx(ctx).Info().Msg("creating node group")
	ng, err := h.provider.CreateNodeGroup(ctx, name, cdef)
	if err != nil {
		return cluster, err
	}

	zerolog.Ctx(ctx).Info().Msg("updating metadata with new nodes")
	var mns []metadata.Node
	cluster.Status = metadata.ClusterConnecting
	if err := h.db.Update(ctx, func(tx *bolt.Tx) error {
		var err error
		tctx := metadata.WithTransactionContext(ctx, tx)
		cluster, err = h.db.UpdateCluster(tctx, cluster)
		if err != nil {
			return err
		}

		mns, err = h.db.CreateNodes(tctx, cluster.ID, ng.Nodes)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return cluster, err
	}

	var ns = make([]p2plab.Node, len(mns))
	for i, n := range mns {
		ns[i] = controlapi.NewNode(h.client, n)
	}

	if err := nodes.WaitHealthy(ctx, ns); err != nil {
		return cluster, err
	}

	zerolog.Ctx(ctx).Info().Msg("updating cluster metadata")
	cluster.Status = metadata.ClusterCreated
	return h.db.UpdateCluster(ctx, cluster)
}
