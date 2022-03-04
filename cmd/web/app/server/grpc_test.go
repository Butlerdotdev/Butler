package server

import (
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	"sync"
	"testing"
)

func TestFailToListen(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	server, err := StartGRPCServer(&GRPCServerParams{
		HostPort: ":-1",
		Logger:   logger,
	})
	assert.Nil(t, server)
	assert.EqualError(t, err, "failed to listen on gRPC port: listen tcp: address -1: invalid port")

}

func TestFailToServe(t *testing.T) {
	ltsn := bufconn.Listen(0)
	ltsn.Close()
	core, logs := observer.New(zap.NewAtomicLevelAt(zapcore.ErrorLevel))
	var wg sync.WaitGroup
	wg.Add(1)

	logger := zap.New(core)
	serveGRPC(grpc.NewServer(), ltsn, &GRPCServerParams{
		Logger: logger,
		OnError: func(e error) {
			assert.Equal(t, 1, len(logs.All()))
			assert.Equal(t, "Could not launch gRPC server", logs.All()[0].Message)
			wg.Done()
		},
	})
	wg.Wait()
}
