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
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net"
	"net/http"
)

// HTTPServerParams is a struct of http parameters
type HTTPServerParams struct {
	HostPort string
	Logger   *zap.Logger
}

// StartHttpServer based on the given parameters
func StartHttpServer(params *HTTPServerParams) (*http.Server, error) {
	params.Logger.Info("Starting Butler HTTP Server", zap.String("http host-port", params.HostPort))

	errorLog, _ := zap.NewStdLogAt(params.Logger, zapcore.ErrorLevel)
	server := &http.Server{
		Addr:     params.HostPort,
		ErrorLog: errorLog,
	}

	// Add TLS

	listener, err := net.Listen("tcp", params.HostPort)
	if err != nil {
		return nil, err
	}

	serveHTTP(server, listener, params)
	return server, nil
}

func serveHTTP(server *http.Server, listener net.Listener, params *HTTPServerParams) {

	go func() {
		var err error
		err = server.Serve(listener)
		if err != nil {
			if err != http.ErrServerClosed {
				params.Logger.Error("Could not start HTTP server", zap.Error(err))
			}
		}
	}()
}
