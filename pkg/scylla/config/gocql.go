// Package config
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
// Temporary file, this will be removed after readers/writers have been implemented. Only using this for the
// utils.bootstrap.go and implementing a keyspace. Keyspace in the future will be initialized via entry script once
// we actually have a Dockerfile / config files for scylla properly setup.

package config

import (
	"github.com/butdotdev/butler/pkg/scylla/db"
	"github.com/scylladb/gocqlx/v2"

	"github.com/gocql/gocql"
	"github.com/spf13/pflag"
	"time"
)

var config = struct {
	DB       gocql.ClusterConfig
	Password gocql.PasswordAuthenticator
}{}

func init() {
	config.DB = *gocql.NewCluster()

	config.DB.Consistency = gocql.LocalOne
	config.DB.PoolConfig.HostSelectionPolicy = gocql.TokenAwareHostPolicy(gocql.RoundRobinHostPolicy())

	pflag.StringArrayVar(
		&config.DB.Hosts,
		"hosts", []string{"127.0.0.1"},
		"cluster nodes address list")
	pflag.DurationVar(
		&config.DB.Timeout,
		"timeout", 60*time.Second,
		"connection timeout")

	pflag.DurationVar(
		&config.DB.ConnectTimeout,
		"dial-timeout", 30*time.Second,
		"initial dial timeout")

	pflag.StringVar(
		&config.Password.Username,
		"username", "",
		"password based authentication username")
	pflag.StringVar(
		&config.Password.Password,
		"password", "",
		"password based authentication password")
}

func Config() gocql.ClusterConfig {

	var tmp = config.DB
	if config.Password.Username != "" {
		tmp.Authenticator = config.Password
	}
	return tmp
}

func Session() (*gocql.Session, error) {
	return gocql.NewSession(Config())
}

func Keyspace() (gocqlx.Session, error) {
	cfg := Config()
	cfg.Keyspace = db.KeySpace
	return gocqlx.WrapSession(gocql.NewSession(cfg))
}
