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
	"fmt"
	"github.com/butdotdev/butler/ports"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
)

type Service struct {
	AdminPort      int
	Logger         *zap.Logger
	signalsChannel chan os.Signal

	Admin *AdminServer
}

func NewService(adminPort int) *Service {
	signalsChannel := make(chan os.Signal, 1)
	//hc

	signal.Notify(signalsChannel, os.Interrupt, syscall.SIGTERM)

	return &Service{
		Admin:          NewAdminServer(ports.PortToHostPort(adminPort)),
		signalsChannel: signalsChannel,
		//hc
	}
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

	s.Admin.initFromViper(v, s.Logger)

	if err := s.Admin.Serve(); err != nil {
		return fmt.Errorf("cannot start the admin server: %w", err)
	}
	return nil
}
