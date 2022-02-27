package app

import (
	"github.com/butdotdev/butler/ports"
	"github.com/spf13/viper"
)
const (
	webHTTPHostPort = "web.http-server.host-port"
)

// WebOptions struct holds config for web
type WebOptions struct {
	// WebHTTPHostPort is the host:port address that the web service listens in on for http requests
	WebHTTPHostPort string
}

func (wOpts *WebOptions) InitFromViper(v *viper.Viper) *WebOptions {
	wOpts.WebHTTPHostPort = ports.FormatHostPort(v.GetString(webHTTPHostPort))

	return wOpts
}