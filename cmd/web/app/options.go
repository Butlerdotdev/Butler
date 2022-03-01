// Package app
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

package app

import "go.uber.org/zap"

type options struct {
	logger *zap.Logger
}

type Option func(c *options)

var Options options

func (options) Logger(logger *zap.Logger) Option {
	return func(b *options) {
		b.logger = logger
	}
}

func (o options) apply(opts ...Option) options {
	ret := options{}

	for _, opt := range opts {
		opt(&ret)
	}
	if ret.logger == nil {
		ret.logger = zap.NewNop()
	}
	return ret
}
