// Package flags
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

package flags

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/butdotdev/butler/pkg/healthcheck"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net"
	"net/http"
)

const (
	HTTPHostPort = "health-server"
)

// Server runs an HTTP server with hc endpoints
type Server struct {
	logger   *zap.Logger
	hostPort string
	hc       *healthcheck.HealthCheck
	mux      *http.ServeMux
	server   *http.Server
}

// NewServer creates a new http server
func NewServer(hostPort string) *Server {
	return &Server{
		hostPort: hostPort,
		logger:   zap.NewNop(),
		hc:       healthcheck.New(),
		mux:      http.NewServeMux(),
	}
}

// HC returns the reference to HeathCheck.
func (s *Server) HC() *healthcheck.HealthCheck {
	return s.hc
}

// setLogger inits a new logger
func (s *Server) setLogger(logger *zap.Logger) {
	s.logger = logger
	s.hc.SetLogger(logger)
}

// AddFlags registers CLI flags.
func (s *Server) AddFlags(flagSet *flag.FlagSet) {
	flagSet.String(HTTPHostPort, s.hostPort, fmt.Sprintf("The host:port (e.g. 127.0.0.1%s or %s) for the admin server, including health check, /healthy, etc.", s.hostPort, s.hostPort))
}

// InitFromViper inits the server with props retrieved from viper
func (s *Server) initFromViper(v *viper.Viper, logger *zap.Logger) {
	s.setLogger(logger)
	s.hostPort = v.GetString(HTTPHostPort)
}

// Handle adds a new hanlder to the HTTP server
func (s *Server) Handle(path string, handler http.Handler) {
	s.mux.Handle(path, handler)
}

// Serve starts a new HTTP Server
func (s *Server) Serve() error {
	l, err := net.Listen("tcp", s.hostPort)
	if err != nil {
		s.logger.Error("server failed to listen", zap.Error(err))
		return err
	}
	s.serveWithListener(l)

	s.logger.Info(
		"server started",
		zap.String("http.host-port", l.Addr().String()),
		zap.Stringer("health-status", s.hc.Get()))
	return nil
}

func (s *Server) serveWithListener(l net.Listener) {
	router := mux.NewRouter()
	router.HandleFunc("/healthy", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	})
	s.logger.Info("Mounting health check on server", zap.String("route", "/healthy"))
	s.mux.Handle("/healthy", s.hc.Handler())
	errorLog, _ := zap.NewStdLogAt(s.logger, zapcore.ErrorLevel)
	s.server = &http.Server{
		ErrorLog: errorLog,
		Handler:  router,
	}

	s.logger.Info("Starting HTTP server", zap.String("http-addr", s.hostPort))
	go func() {
		switch err := s.server.Serve(l); err {
		case nil, http.ErrServerClosed:
			// normal exit, nothing to do
		default:
			s.logger.Error("failed to serve", zap.Error(err))
			s.hc.Set(healthcheck.Broken)
		}
	}()
}

// Close stops the HTTP server
func (s *Server) Close() error {
	return s.server.Shutdown(context.Background())
}
