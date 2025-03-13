// Package docker defines an adapter for Docker.
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

package docker

import (
	"butler/internal/adapters/exec"
	"context"
	"go.uber.org/zap"
)

// DockerAdapter provides a high-level interface for Docker operations.
type DockerAdapter struct {
	client *DockerClient
	logger *zap.Logger
}

// NewDockerAdapter initializes the adapter.
func NewDockerAdapter(execAdapter exec.ExecAdapter, logger *zap.Logger) *DockerAdapter {
	return &DockerAdapter{
		client: NewDockerClient(execAdapter, logger),
		logger: logger,
	}
}

// ExecuteKubectlCommand runs a generic kubectl command with provided arguments.
func (d *DockerAdapter) ExecuteCommand(ctx context.Context, args ...string) (string, error) {
	return d.client.ExecuteCommand(ctx, args...)
}
