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
	"fmt"

	"github.com/butlerdotdev/butler/internal/logger"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// NewBootstrapCmd creates the bootstrap command.
func NewBootstrapCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bootstrap",
		Short: "Bootstrap the Butler management cluster",
		Long:  `Bootstraps the Butler management cluster with the provided configuration. Requires a subcommand to be called specifying the provider.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			log := logger.GetLogger()
			log.Error("Bootstrap needs to be run with a subcommand specifying the provider (e.g., 'butler bootstrap proxmox' or 'butler bootstrap nutanix').")
			return fmt.Errorf("bootstrap needs to be run with a subcommand specifying the provider (e.g., 'butler bootstrap proxmox' or 'butler bootstrap nutanix')")
		},
	}

	// Support CLI-based configuration file override
	cmd.Flags().String("config", "", "Path to configuration file")
	viper.BindPFlag("config", cmd.Flags().Lookup("config"))

	return cmd
}
