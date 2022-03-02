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

import (
	"fmt"
	"net"
	"net/url"
	"runtime"
	"time"

	"go.uber.org/zap"

	"github.com/go-graphite/go-carbon/helper"
	"github.com/go-graphite/go-carbon/points"
)

type statFunc func()

type statModule interface {
	Stat(send helper.StatCallback)
}

type Collector struct {
	helper.Stoppable
	graphPrefix    string
	metricInterval time.Duration
	endpoint       string
	data           chan *points.Points
	stats          []statFunc
	logger         *zap.Logger //nolint:unused,structcheck
}

func RuntimeStat(send helper.StatCallback) {
	send("GOMAXPROCS", float64(runtime.GOMAXPROCS(-1)))
	send("NumGoroutine", float64(runtime.NumGoroutine()))
}

func NewCollector(app *App) *Collector {
	// app locked by caller

	c := &Collector{
		graphPrefix:    app.Config.Common.GraphPrefix,
		metricInterval: app.Config.Common.MetricInterval,
		data:           make(chan *points.Points, 4096),
		endpoint:       app.Config.Common.MetricEndpoint,
		stats:          make([]statFunc, 0),
	}

	c.Start()

	logger := app.Logger

	endpoint, err := url.Parse(c.endpoint)
	if err != nil {
		logger.Error("metric-endpoint parse error", zap.Error(err))
		c.endpoint = "local"
	}

	logger = logger.With(zap.String("endpoint", c.endpoint))
	c.logger = logger

	if c.endpoint == "local" {
		// sender worker
		storeFunc := app.Cache.Add

		c.Go(func(exit chan bool) {
			for {
				select {
				case <-exit:
					return
				case p := <-c.data:
					storeFunc(p)
				}
			}
		})
	} else {
		chunkSize := 32768
		if endpoint.Scheme == "udp" {
			chunkSize = 1000 // nc limitation (1024 for udp) and mtu friendly
		}

		c.Go(func(exit chan bool) {
			points.Glue(exit, c.data, chunkSize, time.Second, func(chunk []byte) {

				var conn net.Conn
				var err error
				defaultTimeout := 5 * time.Second

				// send data to endpoint
			SendLoop:
				for {

					// check exit
					select {
					case <-exit:
						break SendLoop
					default:
						// pass
					}

					// close old broken connection
					if conn != nil {
						conn.Close()
						conn = nil
					}

					conn, err = net.DialTimeout(endpoint.Scheme, endpoint.Host, defaultTimeout)
					if err != nil {
						logger.Error("dial failed", zap.Error(err))
						time.Sleep(time.Second)
						continue SendLoop
					}

					err = conn.SetDeadline(time.Now().Add(defaultTimeout))
					if err != nil {
						logger.Error("conn.SetDeadline failed", zap.Error(err))
						time.Sleep(time.Second)
						continue SendLoop
					}

					_, err := conn.Write(chunk)
					if err != nil {
						logger.Error("conn.Write failed", zap.Error(err))
						time.Sleep(time.Second)
						continue SendLoop
					}

					break SendLoop
				}

				if conn != nil {
					conn.Close()
					conn = nil
				}
			})
		})
	}

	sendCallback := func(moduleName string) func(metric string, value float64) {
		return func(metric string, value float64) {
			key := fmt.Sprintf("%s.%s.%s", c.graphPrefix, moduleName, metric)
			logger.Debug("collect", zap.String("metric", key), zap.Float64("value", value))
			select {
			case c.data <- points.NowPoint(key, value):
				// pass
			default:
				logger.Warn("send queue is full. metric dropped",
					zap.String("key", key), zap.Float64("value", value))
			}
		}
	}

	moduleCallback := func(moduleName string, moduleObj statModule) statFunc {
		return func() {
			moduleObj.Stat(sendCallback(moduleName))
		}
	}

	c.stats = append(c.stats, func() {
		RuntimeStat(sendCallback("runtime"))
	})

	if app.Cache != nil {
		c.stats = append(c.stats, moduleCallback("cache", app.Cache))
	}

	// if app.Carbonserver != nil {
	// 	c.stats = append(c.stats, moduleCallback("carbonserver", app.Carbonserver))
	// }

	if app.Receivers != nil {
		for i := 0; i < len(app.Receivers); i++ {
			c.stats = append(c.stats, moduleCallback(app.Receivers[i].Name, app.Receivers[i]))
		}
	}

	if app.Api != nil {
		c.stats = append(c.stats, moduleCallback("grpc", app.Api))
	}

	// collector worker
	c.Go(func(exit chan bool) {
		ticker := time.NewTicker(c.metricInterval)
		defer ticker.Stop()

		for {
			select {
			case <-exit:
				return
			case <-ticker.C:
				c.collect()
			}
		}
	})

	return c
}

func (c *Collector) collect() {
	c.logger.Info("flushing carbon stats")
	for _, stat := range c.stats {
		stat()
	}
}
