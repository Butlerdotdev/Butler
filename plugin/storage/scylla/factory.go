// Package scylla
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

package scylla

import (
	"flag"
	"github.com/butdotdev/butler/pkg/scylla"
	"github.com/butdotdev/butler/pkg/scylla/config"
	sRuleStore "github.com/butdotdev/butler/plugin/storage/scylla/rulestore"
	"github.com/butdotdev/butler/storage/rulestore"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const (
	primaryStorageConfig = "scylla"
)

type Factory struct {
	logger         *zap.Logger
	Options        *Options
	primarySession scylla.Session
	primaryConfig  config.SessionBuilder
}

// CreateKeySpace implements a NewKeySpace
func (f *Factory) CreateKeySpace() (rulestore.Writer, error) {

	return sRuleStore.NewKeySpace(f.primarySession, f.logger), nil
}

// CreateRuleReader implements nothing and is an unused method as of now. Later it will be used to implement a new Reader
func (f *Factory) CreateRuleReader() {
	panic("implement me")
}

// CreateRuleWriter implements nothing and is an unused method as of now. Later it will be used to implement a new Writer
func (f *Factory) CreateRuleWriter() {
	panic("implement me")
}

// NewFactory returns a new factory using a primary config
func NewFactory() *Factory {
	return &Factory{
		Options: NewOptions(primaryStorageConfig),
	}
}

// AddFlags adds option flags
func (f *Factory) AddFlags(flagSet *flag.FlagSet) {
	f.Options.AddFlags(flagSet)
}

// InitFromViper inits options from for viper
func (f *Factory) InitFromViper(v *viper.Viper, logger *zap.Logger) {
	f.Options.InitFromViper(v)
	f.primaryConfig = f.Options.GetPrimary()
}

// InitFromOptions inits options for the primary config
func (f *Factory) InitFromOptions(o *Options) {
	f.Options = o
	f.primaryConfig = o.GetPrimary()
}

// Initialize implements the Initialize storage interface method
func (f *Factory) Initialize(logger *zap.Logger) error {
	f.logger = logger
	primarySession, err := f.primaryConfig.NewSession(logger)
	if err != nil {
		return err
	}
	f.primarySession = primarySession
	return nil
}
