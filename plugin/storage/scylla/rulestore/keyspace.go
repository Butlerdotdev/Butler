// Package rulestore
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

package rulestore

import (
	"context"
	"github.com/butdotdev/butler/pkg/scylla"
	"github.com/butdotdev/butler/pkg/scylla/db"
	"go.uber.org/zap"
)

type KeySpace struct {
	session scylla.Session
	logger  *zap.Logger
}

func NewKeySpace(session scylla.Session, logger *zap.Logger) *KeySpace {
	return &KeySpace{
		session: session,
		logger:  logger,
	}
}
func (k *KeySpace) writeKeySpace(logger *zap.Logger) error {
	if err := k.session.Query(db.KeySpaceCQL).Exec(); err != nil {
		logger.Fatal("ensure keyspace exists: ", zap.Error(err))
	}
	return nil
}

func (k *KeySpace) WriteKeySpace(ctx context.Context, logger *zap.Logger) error {
	if err := k.writeKeySpace(logger); err != nil {
		return err
	}
	return nil

}

func (k *KeySpace) Close() error {
	k.session.Close()
	return nil
}
