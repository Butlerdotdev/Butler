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
	"github.com/butdotdev/butler/cmd/all-in-one/setupcontext"
	"github.com/butdotdev/butler/cmd/cache/carbon"
	"github.com/butdotdev/butler/cmd/docs"
	"github.com/butdotdev/butler/cmd/flags"
	"github.com/butdotdev/butler/cmd/status"
	"github.com/butdotdev/butler/cmd/web/app"
	"github.com/butdotdev/butler/pkg/config"
	"github.com/butdotdev/butler/ports"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"os"
	"time"
)

func main() {
	setupcontext.SetAllInOne()
	svc := flags.NewService(ports.WebHTTP)
	v := viper.New()
	command := &cobra.Command{
		Use:   "butler-all-in-one",
		Short: "butler all in one",
		Long:  `Full butler distribution, if set. `,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := svc.Start(v); err != nil {
				return err
			}
			logger := svc.Logger

			w := app.New(&app.WebParams{
				ServiceName: "butler-web",
				Logger:      logger,
				HealthCheck: svc.HC(),
			})
			webOpts := new(app.WebOptions).InitFromViper(v)
			if err := w.Start(webOpts); err != nil {
				logger.Fatal("Failed to start the web server", zap.Error(err))
			}

			carbonCache := carbon.NewCarbon(&carbon.Config{
				IngestPort:      2003,
				CarbonQueryPort: 8001,
				Logger:          logger,
				MetricInterval:  60 * time.Second,
				GraphiteHost:    "",
			})
			if err := carbon.Start(carbonCache); err != nil {
				logger.Fatal("Failed to start the carbon process", zap.Error(err))
			}

			svc.RunAndThen(func() {
				if err := w.Close(); err != nil {
					logger.Error("failed to cleanly close the http server", zap.Error(err))
				}
				if err := carbon.Stop(carbonCache); err != nil {
					logger.Error("failed to cleanly close the carbon process", zap.Error(err))
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
	)
	if err := command.Execute(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

}
