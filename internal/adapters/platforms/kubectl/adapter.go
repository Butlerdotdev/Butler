// Package kubectl defines an adapter for executing kubectl commands.
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
// distributed under the Apache License, Version 2.0 (the "License");
// limitations under the License.

package kubectl

import (
	"butler/internal/adapters/exec"
	"context"
	"go.uber.org/zap"
)

// KubectlAdapter provides a high-level interface for interacting with Kubernetes via kubectl.
type KubectlAdapter struct {
	client *KubectlClient
}

// NewKubectlAdapter initializes a new KubectlAdapter.
func NewKubectlAdapter(execAdapter exec.ExecAdapter, logger *zap.Logger) *KubectlAdapter {
	return &KubectlAdapter{client: NewKubectlClient(execAdapter, logger)}
}

// ExecuteKubectlCommand runs a generic kubectl command with provided arguments.
func (a *KubectlAdapter) ExecuteCommand(ctx context.Context, args ...string) (string, error) {
	return a.client.ExecuteCommand(ctx, args...)
}
