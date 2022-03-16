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
// As new storage options are added they needed to be added in here and this needs refactored to account for more than
// just scylla.  The base logic is here, but is still pseudo hard coded for scylla.

package storage

import (
	"flag"
	"fmt"
	"github.com/butdotdev/butler/plugin"
	"github.com/butdotdev/butler/plugin/storage/scylla"
	"github.com/butdotdev/butler/storage"
	"github.com/butdotdev/butler/storage/rulestore"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"log"
)

const (
	scyllaStorageType = "scylla"
)

// AllStorageTypes defines all available storage backends
var AllStorageTypes = []string{scyllaStorageType}

// Factory defines a factory structure
type Factory struct {
	FactoryConfig
	factories map[string]storage.Factory
}

// NewFactory returns the factory value
func NewFactory(config FactoryConfig) (*Factory, error) {
	f := &Factory{FactoryConfig: config}
	f.factories = make(map[string]storage.Factory)
	ff, err := f.getFactoryOfType(scyllaStorageType)
	if err != nil {
		return nil, err
	}
	f.factories[scyllaStorageType] = ff

	return f, nil
}

func (f *Factory) getFactoryOfType(factoryType string) (storage.Factory, error) {
	switch factoryType {
	case scyllaStorageType:
		return scylla.NewFactory(), nil

	default:
		return nil, fmt.Errorf("unknown storage type %s. Valid types are %v", factoryType, AllStorageTypes)
	}
}

// Initialize implements storage.Factory.
func (f *Factory) Initialize(logger *zap.Logger) error {

	for _, factory := range f.factories {
		if err := factory.Initialize(logger); err != nil {
			return err
		}
	}
	return nil
}

// CreateKeySpace implements storage.Factory.
func (f *Factory) CreateKeySpace() (rulestore.Writer, error) {
	var writers []rulestore.Writer
	for _, storageType := range f.RuleWriterTypes {
		factory, ok := f.factories[storageType]
		if !ok {
			return nil, fmt.Errorf("no %s backend registered for span store", storageType)
		}
		writer, err := factory.CreateKeySpace()
		if err != nil {
			return nil, err
		}
		writers = append(writers, writer)
	}
	var ruleWriter rulestore.Writer
	if len(f.RuleWriterTypes) == 1 {
		ruleWriter = writers[0]
	} else {
		log.Panic("panic")
	}
	return ruleWriter, nil
}

// InitFromViper implements plugin.Configurable
func (f *Factory) InitFromViper(v *viper.Viper, logger *zap.Logger) {
	for _, factory := range f.factories {
		if conf, ok := factory.(plugin.Configurable); ok {
			conf.InitFromViper(v, logger)
		}
	}
}

// AddFlags implements plugin.Configurable
func (f *Factory) AddFlags(flagSet *flag.FlagSet) {
	for _, factory := range f.factories {
		if conf, ok := factory.(plugin.Configurable); ok {
			conf.AddFlags(flagSet)
		}
	}
}

// AddPipelineFlags adds the flagset of pipeline flags
func (f *Factory) AddPipelineFlags(flagSet *flag.FlagSet) {
	f.AddFlags(flagSet)
}
