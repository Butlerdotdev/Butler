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
// distributed under the Apache License, Version 2.0 (the "License");
// limitations under the License.

package flux

import (
	"butler/pkg/adapters/exec"
	"context"
	"go.uber.org/zap"
)

// FluxAdapter provides a high-level interface for interacting with Kubernetes via flux.
type FluxAdapter struct {
	client *FluxClient
}

// NewFluxAdapter initializes a new KubectlAdapter.
func NewFluxAdapter(execAdapter exec.ExecAdapter, logger *zap.Logger) *FluxAdapter {
	return &FluxAdapter{client: NewFluxClient(execAdapter, logger)}
}

// ExecuteFluxCommand runs a generic kubectl command with provided arguments.
func (f *FluxAdapter) ExecuteCommand(ctx context.Context, args ...string) (string, error) {
	return f.client.ExecuteCommand(ctx, args...)
}
