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
	bootstrap "butler/internal/handlers/bootstrap/proxmox"
	"butler/internal/logger"
	"context"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// NewProxmoxBootstrapCmd creates the bootstrap command for the Proxmox provider
func NewProxmoxBootstrapCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "proxmox",
		Short: "Bootstrap the Butler management cluster with the Proxmox provider",
		Long: `Bootstraps the Butler management cluster with the provided configuration.
This command provisions the necessary infrastructure in Proxmox and applies cluster configurations.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			log := logger.GetLogger()

			// Initialize the Handler
			handler := bootstrap.NewBootstrapHandler(context.Background(), log)
			if err := handler.HandleProvisionCluster(); err != nil {
				log.Error("Cluster provisioning failed", zap.Error(err))
				return err
			}

			log.Info("Butler bootstrap completed successfully! 🎉")
			return nil
		},
	}

	// Support CLI-based configuration file override
	cmd.Flags().String("config", "", "Path to configuration file")
	viper.BindPFlag("config", cmd.Flags().Lookup("config"))

	return cmd
}
