// Package generate provides utilities for Butler, such as documentation and metadata generation.
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
)

// NewGenerateCmd creates the generate command.
func NewGenerateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate utilities for Butler",
		Run: func(cmd *cobra.Command, args []string) {
			log := logger.GetLogger()
			log.Info("Generate command executed")
		},
	}

	return cmd
}
