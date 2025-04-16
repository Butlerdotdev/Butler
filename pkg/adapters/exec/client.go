// Package exec defines an adapter for executing commands natively within Butler.
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

package exec

import (
	"bytes"
	"context"
	"os/exec"
	"time"

	"github.com/butlerdotdev/butler/pkg/adapters/exec/models"
	"go.uber.org/zap"
)

// Client implements ExecAdapter.
type Client struct {
	logger *zap.Logger
}

// NewClient initializes the Exec Adapter.
func NewClient(logger *zap.Logger) *Client {
	return &Client{logger: logger}
}

func (c *Client) RunCommand(ctx context.Context, cmd string, args ...string) (models.CommandResult, error) {
	c.logger.Info("Executing command", zap.String("command", cmd), zap.Strings("args", args))

	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	command := exec.CommandContext(ctx, cmd, args...)
	var stdout, stderr bytes.Buffer
	command.Stdout = &stdout
	command.Stderr = &stderr

	err := command.Run()

	result := models.CommandResult{
		Command: cmd,
		Args:    args,
		Stdout:  stdout.String(),
		Stderr:  stderr.String(),
		Success: err == nil,
	}

	c.logger.Info("Command execution completed",
		zap.String("stdout", result.Stdout),
		zap.String("stderr", result.Stderr),
		zap.Bool("success", result.Success),
	)

	return result, err
}
