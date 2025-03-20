// Package generate provides utilities for Butler, including documentation generation.
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

package generate

import (
	"butler/internal/logger"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// NewDocsCmd creates the docs generation command.
func NewDocsCmd(rootCmd *cobra.Command) *cobra.Command {
	var outputDir string
	cmd := &cobra.Command{
		Use:   "docs",
		Short: "Generate documentation for Butler",
		Run: func(cmd *cobra.Command, args []string) {
			log := logger.GetLogger()
			outputDir := viper.GetString("docs.outputDir")

			// Generate markdown docs
			err := doc.GenMarkdownTree(rootCmd, outputDir)
			if err != nil {
				log.Error("Failed to generate documentation", zap.Error(err))
				return
			}

			log.Info("Documentation generated successfully!", zap.String("output", outputDir))
		},
	}

	cmd.Flags().StringVarP(&outputDir, "output", "o", "./docs", "Directory to save generated documentation")
	if err := viper.BindPFlag("docs.outputDir", cmd.Flags().Lookup("output")); err != nil {
		logger.GetLogger().Fatal("Failed to bind flag", zap.Error(err))
	}
	return cmd
}
