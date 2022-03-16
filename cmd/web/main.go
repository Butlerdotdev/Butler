// Package main
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

package main

import (
	"fmt"
	"github.com/butdotdev/butler/cmd/docs"
	"github.com/butdotdev/butler/cmd/flags"
	"github.com/butdotdev/butler/cmd/status"
	"github.com/butdotdev/butler/cmd/web/app"
	"github.com/butdotdev/butler/pkg/config"
	"github.com/butdotdev/butler/plugin/storage"
	"github.com/butdotdev/butler/ports"
	"github.com/butdotdev/butler/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"log"
	"os"
)

const serviceName = "butler-web"

func main() {
	svc := flags.NewService(ports.WebHTTP)
	v := viper.New()

	ruleFactory, err := storage.NewFactory(storage.FactoryConfigFromEnvAndCli(os.Args, os.Stderr))
	if err != nil {
		log.Fatalf("cannot initialize rule factory: %v", err)
	}
	command := &cobra.Command{
		Use:   "butler-web",
		Short: "butler web is the main http server for butler",
		Long:  `Butler web is the server that runs and serves the butler frontend.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := svc.Start(v); err != nil {
				return err
			}
			logger := svc.Logger

			logger.Info("Bootstrap database...")
			utils.CreateKeyspace(logger) // need this until I add an actual entrypoint that creates the keyspace.

			ruleFactory.InitFromViper(v, logger)
			if err := ruleFactory.Initialize(logger); err != nil {
				logger.Fatal("Failed to init rules factory", zap.Error(err))
			}

			w := app.New(&app.WebParams{
				ServiceName: serviceName,
				Logger:      logger,
				HealthCheck: svc.HC(),
			})
			webOpts := new(app.WebOptions).InitFromViper(v)
			if err := w.Start(webOpts); err != nil {
				logger.Fatal("Failed to start the web server", zap.Error(err))
			}
			svc.RunAndThen(func() {
				if err := w.Close(); err != nil {
					logger.Error("failed to cleanly close the http server", zap.Error(err))
				}
			})
			return nil

		},
	}

	command.AddCommand(status.Command(v, ports.WebHTTP))
	command.AddCommand(docs.Command(v))

	config.AddFlags(
		v,
		command,
		svc.AddFlags,
		app.AddFlags,
		ruleFactory.AddPipelineFlags,
	)

	if err := command.Execute(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
