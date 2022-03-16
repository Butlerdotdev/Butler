// Package utils
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
// Utils for basic DB entry stuff. Will be getting rid of this after the initial writer / readers are created and
// the base docker image that we plan to use for scylla is setup.
// Some of this stuff is usually done via entry scripts.

package utils

import (
	"context"
	"github.com/butdotdev/butler/pkg/scylla/config"
	"github.com/butdotdev/butler/pkg/scylla/db"
	"github.com/scylladb/gocqlx/v2/migrate"
	"go.uber.org/zap"
)

// CreateKeyspace creates the keyspace in the db  --- this is being replaced. Only using it temporarily
func CreateKeyspace(logger *zap.Logger) {
	ses, err := config.Session()
	if err != nil {
		logger.Fatal("session: ", zap.Error(err))
	}
	defer ses.Close()

	if err := ses.Query(db.KeySpaceCQL).Exec(); err != nil {
		logger.Fatal("ensure keyspace exists: ", zap.Error(err))
	}
}

func MigrateKeyspace(logger *zap.Logger) {
	ses, err := config.Keyspace()
	if err != nil {
		logger.Fatal("session: ", zap.Error(err))
	}
	defer ses.Close()

	if err := migrate.Migrate(context.Background(), ses, "../../pkg/scylla/db/cql"); err != nil {
		logger.Fatal("migrate: ", zap.Error(err))
	}
}

func PrintKeyspaceMetadata(logger *zap.Logger) {
	ses, err := config.Keyspace()
	if err != nil {
		logger.Fatal("session: ", zap.Error(err))
	}
	defer ses.Close()

	m, err := ses.KeyspaceMetadata(db.KeySpace)
	if err != nil {
		logger.Fatal("keyspace metadata: ", zap.Error(err))
	}

	logger.Info("Keyspace Metadata: ", zap.Any("MetaData: ", *m))
}
