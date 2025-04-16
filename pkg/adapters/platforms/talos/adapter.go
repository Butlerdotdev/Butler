// Package talos defines an adapter for Talos and bootstrapping the OS.
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

package talos

import (
	"context"

	"github.com/butlerdotdev/butler/pkg/adapters/exec"
	"go.uber.org/zap"
)

// TalosAdapter provides methods for executing Talos commands.
type TalosAdapter struct {
	client *TalosClient
	logger *zap.Logger
}

// NewTalosAdapter initializes a new Talos adapter.
func NewTalosAdapter(execAdapter exec.ExecAdapter, logger *zap.Logger) *TalosAdapter {
	return &TalosAdapter{
		client: NewTalosClient(execAdapter, logger),
		logger: logger,
	}
}

// ExecuteKubectlCommand runs a generic kubectl command with provided arguments.
func (t *TalosAdapter) ExecuteCommand(ctx context.Context, args ...string) (string, error) {
	return t.client.ExecuteCommand(ctx, args...)
}
