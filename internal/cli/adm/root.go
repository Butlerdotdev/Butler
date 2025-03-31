// Package adm is the root command for the Butler ADM, which serves as the primary entry point butleradm commands.
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

package adm

import (
	"butler/internal/cli/adm/bootstrap"
	"butler/internal/cli/adm/bootstrap/providers"
	"butler/internal/cli/adm/generate"
	"butler/internal/logger"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "butleradm",
	Short: "Butler - Kubernetes as a Service",
	Long: `Butler is an enterprise-grade Kubernetes as a Service Platform.
It supports cluster lifecycle management, infrastructure provisioning, 
and compliance enforcement using a declarative approach.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		cfgFile, _ = cmd.Flags().GetString("config")
		initConfig()
	},
	Run: func(cmd *cobra.Command, args []string) {
		log := logger.GetLogger()
		log.Info("Welcome to Butler CLI! Use --help to view available commands.")
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

func init() {
	// Global configuration flag
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "Path to configuration file")
	viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config"))

	// Initialize configuration
	cobra.OnInitialize(initConfig)

	// Register subcommands
	RegisterCommands()
}

func initConfig() {
	log := logger.GetLogger()

	// If the user specified a config file via --config
	if cfgFile != "" {
		log.Info("Explicit config file detected", zap.String("cfgFile", cfgFile))
		viper.SetConfigFile(cfgFile)

		if err := viper.ReadInConfig(); err != nil {
			log.Fatal("Failed to read config file", zap.Error(err))
		} else {
			log.Info("Using config file", zap.String("file", viper.ConfigFileUsed()))
		}
	} else {
		// Fallback to default config search locations
		viper.SetConfigName("bootstrap")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
		viper.AddConfigPath("$HOME/.butler")

		if err := viper.ReadInConfig(); err != nil {
			log.Warn("No config file found", zap.Error(err))
		} else {
			log.Info("Using default config file", zap.String("file", viper.ConfigFileUsed()))
		}
	}

	viper.AutomaticEnv()
	viper.SetEnvPrefix("BUTLER")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
}

// RegisterCommands explicitly registers all subcommands
func RegisterCommands() {
	bootstrapCmd := bootstrap.NewBootstrapCmd()
	bootstrapCmd.AddCommand(providers.NewNutanixBootstrapCmd(rootCmd))
	bootstrapCmd.AddCommand(providers.NewProxmoxBootstrapCmd())
	rootCmd.AddCommand(bootstrapCmd)

	genCmd := generate.NewGenerateCmd()
	genCmd.AddCommand(generate.NewDocsCmd(rootCmd))
	rootCmd.AddCommand(genCmd)
}

// GetRootCmd returns the root command
func GetRootCmd() *cobra.Command {
	return rootCmd
}
