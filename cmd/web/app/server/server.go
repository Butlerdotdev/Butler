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
	"github.com/butdotdev/butler/cmd/web/app/handler"
	"github.com/butdotdev/butler/pkg/config/tlscfg"
	"github.com/butdotdev/butler/pkg/healthcheck"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net"
	"net/http"
)

// HTTPServerParams is a struct of http parameters
type HTTPServerParams struct {
	HostPort    string
	Logger      *zap.Logger
	HealthCheck *healthcheck.HealthCheck
	TLSConfig   tlscfg.Options
}

// StartHttpServer based on the given parameters
func StartHttpServer(params *HTTPServerParams) (*http.Server, error) {
	router := mux.NewRouter()
	spa := handler.SpaHandler{StaticPath: "../../butler-ui/packages/butler-ui/build", IndexPath: "../../butler-ui/packages/butler-ui/build/index.html"}
	router.PathPrefix("/").Handler(spa)
	params.Logger.Info("Starting the Butler HTTP Server")
	errorLog, _ := zap.NewStdLogAt(params.Logger, zapcore.ErrorLevel)

	server := &http.Server{
		Addr:     params.HostPort,
		ErrorLog: errorLog,
		Handler:  router,
	}
	if params.TLSConfig.Enabled {
		tlsCfg, err := params.TLSConfig.Config(params.Logger) //checks for certs
		if err != nil {
			return nil, err
		}
		server.TLSConfig = tlsCfg
	}
	listener, err := net.Listen("tcp", params.HostPort)
	if err != nil {
		return nil, err
	}
	serveHTTP(server, listener, params)
	params.Logger.Info("Butler HTTP server started", zap.String("http-addr", params.HostPort))
	return server, nil
}

func serveHTTP(server *http.Server, listener net.Listener, params *HTTPServerParams) {
	go func() {
		var err error
		if params.TLSConfig.Enabled {
			err = server.ServeTLS(listener, "", "")
		} else {
			err = server.Serve(listener)
		}
		if err != nil {
			if err != http.ErrServerClosed {
				params.Logger.Error("Could Not Start HTTP Server", zap.Error(err))
			}
		}
		params.HealthCheck.Set(healthcheck.Unavailable)
	}()
}
