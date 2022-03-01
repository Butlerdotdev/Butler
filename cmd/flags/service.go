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
	"flag"
	"fmt"
	"github.com/butdotdev/butler/pkg/healthcheck"
	"github.com/butdotdev/butler/ports"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
)

type Service struct {
	httpPort        int
	Logger          *zap.Logger
	signalsChannel  chan os.Signal
	hcStatusChannel chan healthcheck.Status

	Server *Server
}

func NewService(httpPort int) *Service {
	signalsChannel := make(chan os.Signal, 1)
	hcStatusChannel := make(chan healthcheck.Status)

	signal.Notify(signalsChannel, os.Interrupt, syscall.SIGTERM)

	return &Service{
		Server:          NewServer(ports.PortToHostPort(httpPort)),
		signalsChannel:  signalsChannel,
		hcStatusChannel: hcStatusChannel,
	}
}

// AddFlags registers CLI flags.
func (s *Service) AddFlags(flagSet *flag.FlagSet) {
	AddConfigFileFlag(flagSet)
	if false {
		AddLoggingFlag(flagSet)
	} else {
		s.Server.AddFlags(flagSet)
	}

}

// SetHealthCheckStatus sets status of healthcheck
func (s *Service) SetHealthCheckStatus(status healthcheck.Status) {
	s.hcStatusChannel <- status
}

func (s *Service) Start(v *viper.Viper) error {
	if err := TryLoadConfigFile(v); err != nil {
		return fmt.Errorf("cannot load config file: %w", err)
	}

	sFlags := new(SharedFlags).InitFromViper(v)
	newProdConfig := zap.NewProductionConfig()
	newProdConfig.Sampling = nil
	if logger, err := sFlags.NewLogger(newProdConfig); err == nil {
		s.Logger = logger
	} else {
		return fmt.Errorf("cannot create logger: %w", err)
	}

	s.Server.initFromViper(v, s.Logger)

	if err := s.Server.Serve(); err != nil {
		return fmt.Errorf("cannot start the admin server: %w", err)
	}
	return nil
}

// HC returns the reference to HeathCheck.
func (s *Service) HC() *healthcheck.HealthCheck {
	return s.Server.HC()
}

// RunAndThen sets the health check to Ready and blocks until SIGTERM is received.
// If then runs the shutdown function and exits.
func (s *Service) RunAndThen(shutdown func()) {
	s.HC().Ready()

statusLoop:
	for {
		select {
		case status := <-s.hcStatusChannel:
			s.HC().Set(status)
		case <-s.signalsChannel:
			break statusLoop
		}
	}

	s.Logger.Info("Shutting down")
	s.HC().Set(healthcheck.Unavailable)

	if shutdown != nil {
		shutdown()
	}

	s.Server.Close()
	s.Logger.Info("Shutdown complete")
}
