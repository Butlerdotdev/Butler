// Package flux defines an adapter for executing flux commands.
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

package flux

import (
	"context"
	"fmt"

	"github.com/butlerdotdev/butler/pkg/adapters/exec"

	"go.uber.org/zap"
)

// FluxClient executes kubectl commands.
type FluxClient struct {
	execAdapter exec.ExecAdapter
	logger      *zap.Logger
}

// NewFluxClient initializes a new Kubectl client.
func NewFluxClient(execAdapter exec.ExecAdapter, logger *zap.Logger) *FluxClient {
	return &FluxClient{execAdapter: execAdapter, logger: logger}
}

// ExecuteFluxCommand runs a generic kubectl command with provided arguments.
func (c *FluxClient) ExecuteCommand(ctx context.Context, args ...string) (string, error) {
	result, err := c.execAdapter.RunCommand(ctx, "flux", args...)
	if err != nil {
		return "", fmt.Errorf("flux command failed: %w", err)
	}
	return result.Stdout, nil
}
