// Package carbon
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

package carbon

// Follows format from here
// https://github.com/go-graphite/go-carbon/blob/98c69c8600966ef8b42f57944004dde177b1374c/carbon/app.go

import (
	"github.com/go-graphite/go-carbon/api"
	"github.com/go-graphite/go-carbon/cache"
	"github.com/go-graphite/go-carbon/receiver"
	"go.uber.org/zap"
	"net"
	"sync"
	"time"

	_ "github.com/go-graphite/go-carbon/receiver/tcp"
	_ "github.com/go-graphite/go-carbon/receiver/udp"
)

type NamedReceiver struct {
	receiver.Receiver
	Name string
}

type Options struct {
	Listen     string
	Enabled    bool
	BufferSize int
}

type Common struct {
	GraphPrefix    string
	MetricInterval time.Duration
	MetricEndpoint string
}

type Config struct {
	GRPCAddress   string
	CarbonAddress string
	Logger        *zap.Logger
	Common        *Common
}

// type App struct {
// 	sync.RWMutex
// 	Config         *Config
// 	Api            *api.Api
// 	Cache          *cache.Cache
// 	Receivers      []*NamedReceiver
// 	CarbonLink     *cache.CarbonlinkListener
// 	Persister      *persister.Whisper
// 	Carbonserver   *carbonserver.CarbonserverListener
// 	Tags           *tags.Tags
// 	Collector      *Collector // (!!!) Should be re-created on every change config/modules
// 	PromRegisterer prometheus.Registerer
// 	PromRegistry   *prometheus.Registry
// 	exit           chan bool
// 	FlushTraces    func()
// }

type App struct {
	sync.RWMutex
	Api         *api.Api
	Config      *Config
	Cache       *cache.Cache
	Collector   *Collector
	Logger      *zap.Logger
	exit        chan bool
	FlushTraces func()
	Receivers   []*NamedReceiver
}

func New(config *Config) *App {
	// TODO: Make this configurable
	var duration, _ = time.ParseDuration("1m0s")
	config.Common = &Common{
		GraphPrefix:    "carbon.agents.{host}",
		MetricInterval: duration,
		MetricEndpoint: "local",
	}
	app := &App{
		Config: config,
		Logger: config.Logger,
		exit:   make(chan bool),
	}

	return app
}

func (app *App) stopAll() {
	// Stop all running processes here

	if app.Receivers != nil {
		for i := 0; i < len(app.Receivers); i++ {
			app.Receivers[i].Stop()
			app.Logger.Debug("receiver stopped", zap.String("name", app.Receivers[i].Name))
		}
		app.Receivers = nil
	}

	if app.Api != nil {
		app.Api.Stop()
		app.Api = nil
		app.Logger.Debug("api stopped")
	}

	// if app.Carbonserver != nil {
	// 	carbonserver := app.Carbonserver
	// 	go func() {
	// 		carbonserver.Stop()
	// 		app.Logger.Debug("carbonserver stopped")
	// 	}()
	// 	app.Carbonserver = nil
	// }

	if app.Collector != nil {
		app.Collector.Stop()
		app.Collector = nil
		app.Logger.Debug("collector stopped")
	}

	if app.Cache != nil {
		app.Cache.Stop()
		app.Cache = nil
		app.Logger.Debug("cache stopped")
	}
}

func (app *App) Stop() (err error) {
	app.Lock()
	defer app.Unlock()
	app.stopAll()
	app.Logger.Info("carbon shutdown complete")

	return nil
}

func (app *App) Start() (err error) {
	app.Lock()
	defer app.Unlock()

	defer func() {
		if err != nil {
			app.stopAll()
		}
	}()

	app.Logger.Info("starting carbon")

	// Starts the cache
	core := cache.New()
	//core.SetMaxSize(conf.Cache.MaxSize)

	app.Cache = core

	// Start gRPC API
	var grpcAddr *net.TCPAddr
	grpcAddr, err = net.ResolveTCPAddr("tcp", app.Config.GRPCAddress)
	if err != nil {
		return
	}

	grpcApi := api.New(core)

	if err = grpcApi.Listen(grpcAddr); err != nil {
		return
	}

	app.Api = grpcApi

	// Starts UDP and TCP Receivers
	app.Receivers = make([]*NamedReceiver, 0)
	var rcv receiver.Receiver
	var rcvOptions map[string]interface{}

	var options = &Options{
		Listen:     app.Config.CarbonAddress,
		Enabled:    true,
		BufferSize: 0,
	}

	if rcvOptions, err = receiver.WithProtocol(options, "udp"); err != nil {
		return
	}
	if rcv, err = receiver.New("udp", rcvOptions, core.Add); err != nil {
		return
	}
	app.Receivers = append(app.Receivers, &NamedReceiver{
		Receiver: rcv,
		Name:     "udp",
	})

	if rcvOptions, err = receiver.WithProtocol(options, "tcp"); err != nil {
		return
	}
	if rcv, err = receiver.New("tcp", rcvOptions, core.Add); err != nil {
		return
	}
	app.Receivers = append(app.Receivers, &NamedReceiver{
		Receiver: rcv,
		Name:     "tcp",
	})

	// Starts Stat Collector For Carbon
	app.Collector = NewCollector(app)

	return nil
}

// TODO: Evaluate if this is needed
func (app *App) Loop() {
	app.RLock()
	exitChan := app.exit
	app.RUnlock()

	if exitChan != nil {
		<-app.exit
	}
}

// TODO: add carbonserver
// TODO: add carbonapi
