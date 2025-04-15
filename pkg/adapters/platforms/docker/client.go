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
	"butler/pkg/adapters/exec"
	"context"
	"fmt"
	"go.uber.org/zap"
)

// DockerClient runs Docker commands.
type DockerClient struct {
	execAdapter exec.ExecAdapter
	logger      *zap.Logger
}

// NewDockerClient initializes the client.
func NewDockerClient(execAdapter exec.ExecAdapter, logger *zap.Logger) *DockerClient {
	return &DockerClient{
		execAdapter: execAdapter,
		logger:      logger,
	}
}

// ExecuteKubectlCommand runs a generic kubectl command with provided arguments.
func (c *DockerClient) ExecuteCommand(ctx context.Context, args ...string) (string, error) {
	result, err := c.execAdapter.RunCommand(ctx, "docker", args...)
	if err != nil {
		return "", fmt.Errorf("kubectl command failed: %w", err)
	}
	return result.Stdout, nil
}
