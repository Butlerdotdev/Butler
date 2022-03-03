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
	"github.com/butdotdev/butler/cmd/web/app/server"
	"github.com/butdotdev/butler/pkg/healthcheck"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net/http"
	"time"
)

// Web holds the values for a new web
type Web struct {
	serviceName string
	logger      *zap.Logger
	hCheck      *healthcheck.HealthCheck
	hServer     *http.Server
	grpcServer  *grpc.Server
}

// WebParams to construct a new web
type WebParams struct {
	ServiceName string
	Logger      *zap.Logger
	HealthCheck *healthcheck.HealthCheck
}

// New constructs a new collector component
func New(params *WebParams) *Web {
	return &Web{
		serviceName: params.ServiceName,
		logger:      params.Logger,
		hCheck:      params.HealthCheck,
	}
}

// Start the component and their dependencies
func (w *Web) Start(builderOpts *WebOptions) error {
	httpServer, err := server.StartHttpServer(&server.HTTPServerParams{
		HostPort:    builderOpts.WebHTTPHostPort,
		HealthCheck: w.hCheck,
		Logger:      w.logger,
	})
	if err != nil {
		return fmt.Errorf("could not start the http server %w", err)
	}
	w.hServer = httpServer

	grpcServer, err := server.StartGRPCServer(&server.GRPCServerParams{
		TLSConfig: builderOpts.TLSGRPC,
		HostPort:  builderOpts.WebGRPCHostPort,
		//Handler:                 *handler.GRPCHandler,
		Logger:                  w.logger,
		MaxReceiveMessageLength: builderOpts.WebGRPCMaxReceiveMessageLength,
		MaxConnectionAge:        builderOpts.WebGRPCMaxConnectionAge,
		MaxConnectionAgeGrace:   builderOpts.WebGRPCMaxConnectionAgeGrace,
	})
	if err != nil {
		return fmt.Errorf("could not start gRPC server %w", err)
	}
	w.grpcServer = grpcServer
	return nil
}

// Close the component and all dependencies
func (w *Web) Close() error {
	if w.grpcServer != nil {
		w.grpcServer.GracefulStop()
	}
	if w.hServer != nil {
		timeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		if err := w.hServer.Shutdown(timeout); err != nil {
			w.logger.Fatal("failed to stop the main HTTP server", zap.Error(err))
		}
		defer cancel()
	}
	return nil
}
