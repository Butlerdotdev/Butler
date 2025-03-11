// Package platforms defines generic interface contract and factory.
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

package platforms

import (
	"butler/internal/models"
	"context"
)

// PlatformAdapter defines the interface for platform adapters.
type PlatformAdapter interface {
	Configure(ctx context.Context) error
	GenerateConfig(ctx context.Context, config models.TalosConfig) error
	ApplyConfig(ctx context.Context, node, configDir, role string, insecure bool) error
	BootstrapControlPlane(ctx context.Context, node, configDir string) error
	GetStatus(ctx context.Context) (models.PlatformStatus, error)
}
