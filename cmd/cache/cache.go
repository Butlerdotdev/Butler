// Package main
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

package main

// Follows format from here
// https://github.com/go-graphite/go-carbon/blob/98c69c8600966ef8b42f57944004dde177b1374c/carbon/app.go

import (
	"errors"
	"fmt"
	"net"
	"net/url"
	"os"
	"runtime"
	"strings"
	"sync"
)

type App struct {
	sync.RWMutex
	Config         *Config
	Carbon         *cache.carbon
	exit           chan bool
}

func New(config *config) *App {
	app := &App{
		Config:         config,
		exit:           make(chan bool),
	}

	return app
}

func (app *App) stopAll() {
	// Stop all running caches here
}

func (app *App) Stop() {
	app.Lock()
	defer app.Unlock()
	app.stopAll()
}

func (app *App) Start(version string) (err error) {
	app.Lock()
	defer app.Unlock()

	defer func() {
		if err != nil {
			app.stopAll()
		}
	}()
}

func (app *App) Loop() {
	app.RLock()
	exitChan := app.exit
	app.RUnlock()

	if exitChan != nil {
		<-app.exit
	}
}
