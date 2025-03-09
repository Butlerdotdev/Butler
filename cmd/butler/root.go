// Package butler is the root command for the Butler CLI, which serves as the primary entry point.
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

package butler

import (
	"butler/cmd/butler/bootstrap"
	"butler/cmd/butler/generate"
	"butler/internal/logger"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var rootCmd = &cobra.Command{
	Use:   "butler",
	Short: "Butler - Kubernetes as a Service",
	Run: func(cmd *cobra.Command, args []string) {
		logger.GetLogger().Info("Welcome to Butler CLI! Use --help to view available commands.")
	},
}

// Execute runs the CLI
func Execute() {
	logger.InitLogger()
	log := logger.GetLogger()

	if err := rootCmd.Execute(); err != nil {
		log.Fatal("Execution failed", zap.Error(err))
		os.Exit(1)
	}
}

// RegisterCommands explicitly registers all subcommands
func RegisterCommands() {
	rootCmd.AddCommand(bootstrap.NewBootstrapCmd())
	genCmd := generate.NewGenerateCmd()
	genCmd.AddCommand(generate.NewDocsCmd(rootCmd))
	rootCmd.AddCommand(genCmd)
}

func init() {
	cobra.OnInitialize(initConfig)
	RegisterCommands()
}

func initConfig() {
	log := logger.GetLogger()
	viper.SetConfigName("bootstrap")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		log.Info("Using config file", zap.String("file", viper.ConfigFileUsed()))
	} else {
		log.Warn("No config file found, using defaults")
	}
}

// GetRootCmd returns the root command
func GetRootCmd() *cobra.Command {
	return rootCmd
}
