// Package main is the entry point for butleradm.
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

package main

import (
	"github.com/butlerdotdev/butler/internal/cli/adm"
	_ "github.com/butlerdotdev/butler/internal/cli/adm/generate"
	"github.com/butlerdotdev/butler/internal/logger"
)

// main initializes the logger and executes the Butler CLI.
func main() {
	logger.InitLogger()
	adm.Execute()
}
