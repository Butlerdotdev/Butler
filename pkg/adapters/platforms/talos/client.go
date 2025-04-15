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
	"butler/pkg/adapters/exec"
	"context"
	"fmt"
	"go.uber.org/zap"
)

// TalosClient interacts with Talos using the ExecAdapter.
type TalosClient struct {
	execAdapter exec.ExecAdapter
	logger      *zap.Logger
}

// NewTalosClient creates a new Talos client.
func NewTalosClient(execAdapter exec.ExecAdapter, logger *zap.Logger) *TalosClient {
	return &TalosClient{execAdapter: execAdapter, logger: logger}
}

// ExecuteCommand runs a Talos command with provided arguments.
func (c *TalosClient) ExecuteCommand(ctx context.Context, args ...string) (string, error) {
	c.logger.Info("Executing Talos command", zap.Strings("args", args))

	result, err := c.execAdapter.RunCommand(ctx, "talosctl", args...)
	if err != nil {
		return "", fmt.Errorf("talosctl command failed: %w", err)
	}
	return result.Stdout, nil
}
