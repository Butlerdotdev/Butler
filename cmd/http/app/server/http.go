package server

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net"
	"net/http"
)

// HTTPServerParams is a struct of http parameters
type HTTPServerParams struct {
	HostPort       string
	Logger         *zap.Logger
}


// StartHttpServer based on the given parameters
func StartHttpServer(params *HTTPServerParams) (*http.Server, error) {
	params.Logger.Info("Starting Butler HTTP Server", zap.String("http host-port", params.HostPort))

	errorLog, _ := zap.NewStdLogAt(params.Logger, zapcore.ErrorLevel)
	server := &http.Server{
		Addr: params.HostPort,
		ErrorLog: errorLog,
	}

	// Add TLS

	listener, err := net.Listen("tcp", params.HostPort)
	if err != nil {
		return nil, err
	}

	serveHTTP(server, listener, params)
	return server, nil
}

func serveHTTP(server *http.Server, listener net.Listener, params *HTTPServerParams) {

	go func() {
		var err error
		err = server.Serve(listener)
		if err != nil {
			if err != http.ErrServerClosed {
				params.Logger.Error("Could not start HTTP server", zap.Error(err))
			}
		}
	}()
}