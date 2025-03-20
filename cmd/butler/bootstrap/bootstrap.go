// Package bootstrap provides functionality to bootstrap the Butler management cluster.
//
// Copyright (c) 2025, The Butler Authors
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

package bootstrap

import (
	"butler/internal/handlers/bootstrap"
	"butler/internal/logger"
	"context"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// NewBootstrapCmd creates the bootstrap command.
func NewBootstrapCmd() *cobra.Command {
	var configFile string
	cmd := &cobra.Command{
		Use:   "bootstrap",
		Short: "Bootstrap the Butler management cluster",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			log := logger.GetLogger()
			log.Info("Starting Butler management cluster bootstrap...")
			if err := viper.BindPFlag("config", cmd.Flags().Lookup("config")); err != nil {
				log.Fatal("Failed to bind config flag", zap.Error(err))
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			log := logger.GetLogger()

			// Ensure config is loaded
			if err := viper.ReadInConfig(); err != nil {
				log.Fatal("Failed to load config", zap.Error(err))
			}

			// Initialize Handler and Start Bootstrap
			handler := bootstrap.NewBootstrapHandler(context.Background(), log)
			if err := handler.HandleProvisionCluster(); err != nil {
				log.Fatal("Cluster provisioning failed", zap.Error(err))
			}

			log.Info("Butler bootstrap completed successfully! 🎉")
		},
	}
	cmd.Flags().StringVarP(&configFile, "config", "c", "bootstrap.yaml", "Path to bootstrap configuration file")
	return cmd
}
