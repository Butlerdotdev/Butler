// Package server
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

package server

import (
	"fmt"
	"github.com/butdotdev/butler/cmd/web/app/handler"
	"github.com/butdotdev/butler/pkg/config/tlscfg"
	"github.com/butdotdev/butler/pkg/healthcheck"
	"github.com/butdotdev/butler/proto-gen/api_v1"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
	"net"
	"time"
)

// GRPCServerParams is a struct that holds the parameters for a gRPC server
type GRPCServerParams struct {
	TLSConfig               tlscfg.Options
	HostPort                string
	Handler                 *handler.GRPCHandler
	Logger                  *zap.Logger
	OnError                 func(error)
	MaxReceiveMessageLength int
	unavailableChannel      chan healthcheck.Status
	MaxConnectionAge        time.Duration
	MaxConnectionAgeGrace   time.Duration
	HostPortActual          string // set by the server to determine the actual host:port of the server itself
}

// StartGRPCServer based on the parameters
func StartGRPCServer(params *GRPCServerParams) (*grpc.Server, error) {
	var server *grpc.Server
	var grpcOpts []grpc.ServerOption

	if params.MaxReceiveMessageLength > 0 {
		grpcOpts = append(grpcOpts, grpc.MaxRecvMsgSize(params.MaxReceiveMessageLength))
	}
	grpcOpts = append(grpcOpts, grpc.KeepaliveParams(keepalive.ServerParameters{
		MaxConnectionAge:      params.MaxConnectionAge,
		MaxConnectionAgeGrace: params.MaxConnectionAgeGrace,
	}))
	if params.TLSConfig.Enabled {
		tlsCfg, err := params.TLSConfig.Config(params.Logger)
		if err != nil {
			return nil, err
		}
		creds := credentials.NewTLS(tlsCfg)
		grpcOpts = append(grpcOpts, grpc.Creds(creds))
	}

	server = grpc.NewServer(grpcOpts...)
	reflection.Register(server)

	listener, err := net.Listen("tcp", params.HostPort)
	if err != nil {
		return nil, fmt.Errorf("failed to listen on gRPC port: %w", err)
	}
	params.HostPortActual = listener.Addr().String()

	if err := serveGRPC(server, listener, params); err != nil {
		return nil, err
	}
	return server, nil
}

// serveGRPC serves the rpc server
func serveGRPC(server *grpc.Server, listener net.Listener, params *GRPCServerParams) error {

	nGRPCHandler := handler.NewGRPCHandler(params.Logger)
	api_v1.RegisterWebServiceServer(server, nGRPCHandler)
	params.Logger.Info("Starting Butler gRPC server", zap.String("grpc.host-port", params.HostPortActual))
	go func() {
		if err := server.Serve(listener); err != nil {
			params.Logger.Error("Could not launch gRPC server", zap.Error(err))
			if params.OnError != nil {
				params.OnError(err)
			}
			params.unavailableChannel <- healthcheck.Unavailable
		}
	}()
	params.Logger.Info("Butler gRPC server started", zap.String("grpc.host-port", params.HostPortActual))
	return nil
}
