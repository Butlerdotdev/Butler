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
	"context"

	"github.com/butlerdotdev/butler/pkg/adapters/exec/models"
)

// ExecAdapter defines an interface for executing system commands.
type ExecAdapter interface {
	RunCommand(ctx context.Context, cmd string, args ...string) (models.CommandResult, error)
}
