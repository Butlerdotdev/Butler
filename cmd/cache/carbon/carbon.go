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
	"bufio"
	"encoding/json"
	"errors"
	"expvar"
	"fmt"
	"math"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/go-graphite/carbonzipper/mstats"
	pb "github.com/go-graphite/protocol/carbonapi_v3_pb"

	"github.com/peterbourgon/g2g"

	"github.com/go-graphite/carbonmem/mwhisper"
)

var BuildVersion = "(development build)"
var logger *zap.Logger

var Metrics = struct {
	FindRequests  *expvar.Int
	FetchRequests *expvar.Int
}{
	FindRequests:  expvar.NewInt("findRequests"),
	FetchRequests: expvar.NewInt("fetchRequests"),
}

type Config struct {
	IngestPort      int
	CarbonQueryPort int
	Logger          *zap.Logger
	MetricInterval  time.Duration
	GraphiteHost    string
}

type App struct {
	sync.RWMutex
	Config *Config
	Logger *zap.Logger
}

func NewCarbon(config *Config) *App {
	logger = config.Logger

	app := &App{
		Config: config,
		Logger: config.Logger,
	}

	return app
}

func Start(app *App) (err error) {
	Whispers.metrics = make(map[string]*mwhisper.Whisper)

	expvar.NewString("BuildVersion").Set(BuildVersion)
	app.Logger.Info("starting carbonmem")

	expvar.Publish("Whispers", expvar.Func(func() interface{} {
		m := make(map[string]int)
		Whispers.RLock()
		for k, v := range Whispers.metrics {
			m[k] = v.Len()
		}
		Whispers.RUnlock()
		return m
	}))

	if Whispers.epoch0 == 0 {
		Whispers.epoch0 = int(time.Now().Unix())
	}

	if envhost := os.Getenv("GRAPHITEHOST") + ":" + os.Getenv("GRAPHITEPORT"); envhost != ":" || app.Config.GraphiteHost != "" {

		var host string

		switch {
		case envhost != ":" && app.Config.GraphiteHost != "":
			host = app.Config.GraphiteHost
		case envhost != ":":
			host = envhost
		case app.Config.GraphiteHost != "":
			host = app.Config.GraphiteHost
		}

		app.Logger.Info("Using graphite host", zap.String("host", host))

		// register our metrics with graphite
		graphite := g2g.NewGraphite(host, app.Config.MetricInterval, 10*time.Second)

		hostname, _ := os.Hostname()
		hostname = strings.Replace(hostname, ".", "_", -1)

		graphite.Register(fmt.Sprintf("carbon.mem.%s.find_requests", hostname), Metrics.FindRequests)
		graphite.Register(fmt.Sprintf("carbon.mem.%s.fetch_requests", hostname), Metrics.FetchRequests)

		go mstats.Start(app.Config.MetricInterval)

		graphite.Register(fmt.Sprintf("carbon.mem.%s.alloc", hostname), &mstats.Alloc)
		graphite.Register(fmt.Sprintf("carbon.mem.%s.total_alloc", hostname), &mstats.TotalAlloc)
		graphite.Register(fmt.Sprintf("carbon.mem.%s.num_gc", hostname), &mstats.NumGC)
		graphite.Register(fmt.Sprintf("carbon.mem.%s.pause_ns", hostname), &mstats.PauseNS)
	}

	go graphiteServer(app.Config.IngestPort)

	http.HandleFunc("/metrics/find/", accessHandler(false, findHandler)) // TODO: false here refers to debug logging, make configurable?
	http.HandleFunc("/render/", accessHandler(false, renderHandler))

	app.Logger.Info("carbon query http server starting on port", zap.Int("carbon query port", app.Config.CarbonQueryPort))
	go func() {
		http.ListenAndServe(":"+strconv.Itoa(app.Config.CarbonQueryPort), nil)
	}()

	return nil
}

func Stop(app *App) (err error) {
	logger.Info("stopping carbon")
	return nil
}

func parseTopK(query string) (string, int32, bool) {

	// prefix.blah.*.TopK.10m  => "prefix.blah.*", 600, true

	var idx int
	if idx = strings.Index(query, ".TopK."); idx == -1 {
		// not found
		return "", 0, false
	}

	prefix := query[:idx]

	timeIdx := idx + len(".TopK.")

	// look for number followed by 'm' or 's'
	unitsIdx := timeIdx
	for unitsIdx < len(query) && '0' <= query[unitsIdx] && query[unitsIdx] <= '9' {
		unitsIdx++
	}

	// ran off the end or no numbers present
	if unitsIdx == len(query) || unitsIdx == timeIdx {
		return "", 0, false
	}

	multiplier := 0
	switch query[unitsIdx] {
	case 's':
		multiplier = 1
	case 'm':
		multiplier = 60
	default:
		// unknown units
		return "", 0, false
	}

	if unitsIdx != len(query)-1 {
		return "", 0, false
	}

	timeUnits, err := strconv.Atoi(query[timeIdx:unitsIdx])
	if err != nil {
		return "", 0, false
	}

	return prefix, int32(timeUnits * multiplier), true
}

func findHandler(w http.ResponseWriter, req *http.Request) {

	Metrics.FindRequests.Add(1)

	query := req.FormValue("query")
	format := req.FormValue("format")

	if format != "json" && format != "protobuf3" && format != "protobuf" {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	matches := findMetrics(query)

	response := pb.GlobResponse{
		Name:    query,
		Matches: matches,
	}

	var b []byte
	switch format {
	case "json":
		w.Header().Set("Content-Type", "application/json")
		b, _ = json.Marshal(response)
	case "protobuf", "protobuf3":
		w.Header().Set("Content-Type", "application/protobuf")
		b, _ = response.Marshal()
	}
	w.Write(b)
}

func hasMetaCharacters(query string) bool {
	return strings.IndexByte(query, '*') != -1 || strings.IndexByte(query, '[') != -1 || strings.IndexByte(query, '?') != -1
}

func findMetrics(query string) []pb.GlobMatch {
	var topk string

	var globs []mwhisper.Glob

	if strings.Count(query, ".") < Whispers.prefix {
		globs = Whispers.Glob(query)
	} else {
		if m := Whispers.Fetch(query); m != nil {
			if prefix, seconds, ok := parseTopK(query); ok {
				topk = query[len(prefix):]
				if hasMetaCharacters(query) {
					globs = m.TopK(prefix, seconds)
				} else {
					globs = m.Find(prefix)
				}
			} else {
				globs = m.Find(query)
			}
		}
	}

	var matches []pb.GlobMatch
	paths := make(map[string]struct{}, len(globs))
	for _, g := range globs {
		// fix up metric name
		metric := g.Metric + topk
		if _, ok := paths[metric]; !ok {
			m := pb.GlobMatch{
				Path:   metric,
				IsLeaf: g.IsLeaf,
			}
			matches = append(matches, m)
			paths[metric] = struct{}{}
		}
	}

	return matches
}

func renderHandler(w http.ResponseWriter, req *http.Request) {

	Metrics.FetchRequests.Add(1)

	target := req.FormValue("target")
	format := req.FormValue("format")
	from := req.FormValue("from")
	until := req.FormValue("until")

	frint, _ := strconv.Atoi(from)
	unint, _ := strconv.Atoi(until)

	if unint < frint {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	if format != "json" && format != "protobuf3" && format != "protobuf" {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	matches := findMetrics(target)

	var multi pb.MultiFetchResponse

	for _, m := range matches {

		target := m.GetPath()

		var metric string
		if prefix, _, ok := parseTopK(target); ok {
			metric = prefix
		} else {
			metric = target
		}

		metrics := Whispers.Fetch(metric)
		if metrics == nil {
			continue
		}
		points := metrics.Fetch(metric, int32(frint), int32(unint))

		if points == nil {
			continue
		}

		fromTime := points.From
		untilTime := points.Until
		step := points.Step
		response := pb.FetchResponse{
			Name:      target,
			StartTime: int64(fromTime),
			StopTime:  int64(untilTime),
			StepTime:  int64(step),
			Values:    make([]float64, len(points.Values)),
			//IsAbsent:  make([]bool, len(points.Values)),
		}
		// TODO: Address the commented out sections here. Was put this way just to get it to compile. Functionality unknown
		for i, p := range points.Values {
			if math.IsNaN(p) {
				response.Values[i] = p // This was 0
				//response.IsAbsent[i] = true
			} else {
				response.Values[i] = p
				//response.IsAbsent[i] = false
			}
		}

		multi.Metrics = append(multi.Metrics, response)
	}

	var b []byte
	switch format {
	case "json":
		w.Header().Set("Content-Type", "application/json")
		b, _ = json.Marshal(multi)
	case "protobuf", "protobuf3":
		w.Header().Set("Content-Type", "application/protobuf")
		b, _ = multi.Marshal()
	}
	w.Write(b)
}

func graphiteServer(port int) {

	ln, e := net.Listen("tcp", ":"+strconv.Itoa(port))

	if e != nil {
		logger.Error("listen error", zap.Error(e))
	}

	logger.Info("graphite server starting on port", zap.Int("port", port))

	for {
		conn, err := ln.Accept()
		if err != nil {
			logger.Error("error", zap.Error(err))
			continue
		}
		go func(c net.Conn) {
			scanner := bufio.NewScanner(c)
			for scanner.Scan() {
				metric, count, epoch, err := parseGraphite(scanner.Bytes())
				if err != nil {
					continue
				}

				metrics := Whispers.FetchOrCreate(metric)

				metrics.Set(int32(epoch), metric, uint64(count))
			}
			if err := scanner.Err(); err != nil {
				logger.Warn("graphite server: error during scan", zap.Error(err))
			}
			c.Close()
		}(conn)
	}
}

func isspace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\r' || c == '\n'
}

var errParseError = errors.New("graphite: parse error")

func token(b []byte) ([]byte, []byte) {

	if len(b) == 0 {
		return nil, nil
	}

	// munch space
	for len(b) > 0 && isspace(b[0]) {
		b = b[1:]
	}

	var i int
	for i < len(b) && !isspace(b[i]) {
		i++
	}

	return b[:i], b[i:]
}

func parseGraphite(b []byte) (metric string, count int, epoch int, err error) {

	var tok []byte

	tok, b = token(b)
	if len(tok) == 0 {
		return "", 0, 0, errParseError
	}

	metric = string(tok)

	tok, b = token(b)
	if len(tok) == 0 {
		return "", 0, 0, errParseError
	}

	count, err = strconv.Atoi(string(tok))
	if err != nil {
		return "", 0, 0, errParseError
	}

	tok, b = token(b)
	if len(tok) == 0 {
		return "", 0, 0, errParseError
	}

	epoch, err = strconv.Atoi(string(tok))
	if err != nil {
		return "", 0, 0, errParseError
	}

	// check for extra stuff
	tok, b = token(b)
	if len(tok) != 0 || len(b) != 0 {
		return "", 0, 0, errParseError
	}

	return metric, count, epoch, nil
}

type whispers struct {
	sync.RWMutex
	metrics map[string]*mwhisper.Whisper

	windowSize int
	epochSize  int
	epoch0     int
	prefix     int
}

var Whispers whispers

func findNodePrefix(prefix int, metric string) string {

	var found int
	for i, c := range metric {
		if c == '.' {
			found++
			if found >= prefix {
				return metric[:i]
			}
		}
	}
	return metric
}

func (w *whispers) FetchOrCreate(metric string) *mwhisper.Whisper {

	m := w.Fetch(metric)

	if m == nil {
		prefix := findNodePrefix(w.prefix, metric)
		var ok bool
		w.Lock()
		m, ok = w.metrics[prefix]
		if !ok {
			m = mwhisper.NewWhisper(int32(w.epoch0), w.epochSize, w.windowSize, mwhisper.TrigramCutoff(100000))
			w.metrics[prefix] = m
		}
		w.Unlock()
	}

	return m
}

func (w *whispers) Fetch(metric string) *mwhisper.Whisper {
	prefix := findNodePrefix(w.prefix, metric)

	w.RLock()
	m := w.metrics[prefix]
	w.RUnlock()

	return m
}

func (w *whispers) Glob(query string) []mwhisper.Glob {

	query = strings.Replace(query, ".", "/", -1)
	slashes := strings.Count(query, "/")

	w.RLock()
	var glob []mwhisper.Glob
	for m := range w.metrics {
		qm := strings.Replace(m, ".", "/", slashes)
		if trim := strings.Index(qm, "."); trim != -1 {
			qm = qm[:trim]
			m = m[:trim]
		}
		if match, err := filepath.Match(query, qm); err == nil && match {
			glob = append(glob, mwhisper.Glob{Metric: m})
		}
	}

	w.RUnlock()

	return glob
}

func accessHandler(verbose bool, handler http.HandlerFunc) http.HandlerFunc {
	if !verbose {
		return handler
	}
	return func(w http.ResponseWriter, r *http.Request) {
		t0 := time.Now()
		handler(w, r)
		since := time.Since(t0)
		logger.Info(r.RequestURI, zap.Int64("request time (ms)", since.Nanoseconds()/int64(time.Millisecond)))
	}
}
