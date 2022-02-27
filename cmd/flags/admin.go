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
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net"
	"net/http"
)

const (
	adminHTTPHostPort = "admin.http.host-port"
)

type AdminServer struct {
	logger        *zap.Logger
	adminHostPort string
	//hc
	mux    *http.ServeMux
	server *http.Server
}

func NewAdminServer(hostPort string) *AdminServer {
	return &AdminServer{
		adminHostPort: hostPort,
		logger:        zap.NewNop(),
		//hc
		mux: http.NewServeMux(),
	}
}

func (s *AdminServer) setLogger(logger *zap.Logger) {
	s.logger = logger
	//hc
}

func (s *AdminServer) initFromViper(v *viper.Viper, logger *zap.Logger) {
	s.setLogger(logger)
	s.adminHostPort = v.GetString(adminHTTPHostPort)
}

func (s *AdminServer) Handle(path string, handler http.Handler) {
	s.mux.Handle(path, handler)
}

func (s *AdminServer) Serve() error {
	l, err := net.Listen("tcp", s.adminHostPort)
	if err != nil {
		s.logger.Error("Admin server failed to listen", zap.Error(err))
		return err
	}
	s.serveWithListener(l)

	s.logger.Info(
		"Admin server started",
		zap.String("http.host-port", l.Addr().String()),
	)
	return nil
}

func (s *AdminServer) serveWithListener(l net.Listener) {
	s.logger.Info("Mounting health check on admin server", zap.String("route", "/"))
	//hc
	//version.RegisterHandler(s.mux, s.logger)
	errorLog, _ := zap.NewStdLogAt(s.logger, zapcore.ErrorLevel)
	s.server = &http.Server{
		ErrorLog: errorLog,
	}

	s.logger.Info("Starting admin HTTP server", zap.String("http-addr", s.adminHostPort))
	go func() {
		switch err := s.server.Serve(l); err {
		case nil, http.ErrServerClosed:
			// normal exit, nothing to do
		default:
			s.logger.Error("failed to serve", zap.Error(err))
			//hc
		}
	}()
}

func (s *AdminServer) Close() error {
	return s.server.Shutdown(context.Background())

}
