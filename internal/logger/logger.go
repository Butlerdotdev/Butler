// Package logger provides a singleton implementation of a structured logger using zap.
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

package logger

import (
	"go.uber.org/zap"
	"sync"
)

var (
	log  *zap.Logger
	once sync.Once
)

// InitLogger initializes the logger exactly once
func InitLogger() {
	once.Do(func() {
		var err error
		log, err = zap.NewProduction()
		if err != nil {
			panic("failed to initialize zap logger: " + err.Error())
		}
	})
}

// GetLogger returns the zap logger instance
func GetLogger() *zap.Logger {
	if log == nil {
		InitLogger()
	}
	return log
}

// CloseLogger flushes buffered logs
func CloseLogger() {
	if log != nil {
		_ = log.Sync()
	}
}
