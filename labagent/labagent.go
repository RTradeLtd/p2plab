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

package labagent

import (
	"context"
	"io"
	"os"
	"path/filepath"

	"github.com/Netflix/p2plab/daemon"
	"github.com/Netflix/p2plab/daemon/healthcheckrouter"
	"github.com/Netflix/p2plab/downloaders"
	"github.com/Netflix/p2plab/labagent/agentrouter"
	"github.com/Netflix/p2plab/labagent/supervisor"
	"github.com/Netflix/p2plab/pkg/httputil"
	"github.com/rs/zerolog"
)

type LabAgent struct {
	daemon  *daemon.Daemon
	closers []io.Closer
}

func New(root, addr, appRoot, appAddr string, logger *zerolog.Logger, opts ...LabagentOption) (*LabAgent, error) {
	var settings LabagentSettings
	for _, opt := range opts {
		err := opt(&settings)
		if err != nil {
			return nil, err
		}
	}
	os.MkdirAll(root, 0711)
	client, err := httputil.NewClient(httputil.NewHTTPClient(), httputil.WithLogger(logger))
	if err != nil {
		return nil, err
	}

	settings.DownloaderSettings.Client = client
	fs := downloaders.New(filepath.Join(root, "downloaders"), settings.DownloaderSettings)

	s, err := supervisor.New(filepath.Join(root, "supervisor"), appRoot, appAddr, client, fs)
	if err != nil {
		return nil, err
	}

	var closers []io.Closer
	daemon, err := daemon.New("labagent", addr, logger,
		healthcheckrouter.New(),
		agentrouter.New(appAddr, s),
	)
	if err != nil {
		return nil, err
	}
	closers = append(closers, daemon)

	return &LabAgent{
		daemon:  daemon,
		closers: closers,
	}, nil
}

func (a *LabAgent) Close() error {
	for _, closer := range a.closers {
		err := closer.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *LabAgent) Serve(ctx context.Context) error {
	return a.daemon.Serve(ctx)
}
