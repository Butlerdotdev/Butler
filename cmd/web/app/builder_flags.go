// Package app
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

package app

import (
	"flag"
	"github.com/butdotdev/butler/ports"
	"github.com/spf13/viper"
)

const (
	webHTTPHostPort = "non-server"
)

// WebOptions struct holds config for web
type WebOptions struct {
	// WebHTTPHostPort is the host:port address that the web service listens in on for http requests
	WebHTTPHostPort string
}

func AddFlags(flags *flag.FlagSet) {
	flags.String(webHTTPHostPort, ports.PortToHostPort(ports.WebHTTP), "The host:port (e.g. 127.0.0.1:3000 ) of Butler's HTTP server")
}

func (wOpts *WebOptions) InitFromViper(v *viper.Viper) *WebOptions {
	wOpts.WebHTTPHostPort = ports.FormatHostPort(v.GetString(webHTTPHostPort))

	return wOpts
}
