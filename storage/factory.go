// Package storage
// Copyright (c) 2022, The Butler Authors
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
//

package storage

import (
	"github.com/butdotdev/butler/storage/rulestore"
	"go.uber.org/zap"
)

// Factory defines an interface for the rules factory. It currently holds some placeholder methods that can
// be used later to implement the scylla reader/writer
type Factory interface {
	Initialize(logger *zap.Logger) error
	CreateRuleReader()
	CreateRuleWriter()
	CreateKeySpace() (rulestore.Writer, error)
}