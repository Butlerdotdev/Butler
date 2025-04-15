// Package helm defines an adapter for executing helm commands.
//
// Copyright (c) 2025, The Butler Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package helm

import (
	"butler/pkg/adapters/exec"
	"context"

	"go.uber.org/zap"
)

// HelmAdapter provides a high-level interface for interacting with Helm CLI.
type HelmAdapter struct {
	client *HelmClient
}

// NewHelmAdapter initializes a new HelmAdapter.
func NewHelmAdapter(execAdapter exec.ExecAdapter, logger *zap.Logger) *HelmAdapter {
	return &HelmAdapter{client: NewHelmClient(execAdapter, logger)}
}

// ExecuteCommand runs a generic helm command with provided arguments.
func (a *HelmAdapter) ExecuteCommand(ctx context.Context, args ...string) (string, error) {
	return a.client.ExecuteCommand(ctx, args...)
}
