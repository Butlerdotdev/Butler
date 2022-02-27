// Package app
// Copyright (c) 2022, The Butler Authors
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
//

package app

import (
	"context"
	"fmt"
	"github.com/butdotdev/butler/cmd/http/app/server"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type Web struct {
	logger      *zap.Logger
	serviceName string

	hServer *http.Server
}

type WebParams struct {
	ServiceName string
	Logger      *zap.Logger
}

func New(params *WebParams) *Web {
	return &Web{
		serviceName: params.ServiceName,
		logger:      params.Logger,
	}
}

func (w *Web) Start(builderOpts *WebOptions) error {
	httpServer, err := server.StartHttpServer(&server.HTTPServerParams{
		HostPort: builderOpts.WebHTTPHostPort,
		Logger:   w.logger,
	})
	if err != nil {
		return fmt.Errorf("could not start the HTTP server %w", err)
	}
	w.hServer = httpServer
	//w.publishOpts(builderOpts)
	return nil
}

func (w *Web) Close() error {
	if w.hServer != nil {
		timeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		if err := w.hServer.Shutdown(timeout); err != nil {
			w.logger.Fatal("failed to stop the main HTTP server", zap.Error(err))
		}
		defer cancel()
	}
	return nil
}
