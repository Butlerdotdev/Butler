// Package handler
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

package handler

import (
	"context"
	"github.com/butdotdev/butler/proto-gen/api_v1"
	"go.uber.org/zap"
)

// GRPCHandler is the handler for gRPC
type GRPCHandler struct {
	logger                               *zap.Logger
	api_v1.UnimplementedWebServiceServer // this is dumb
}

// NewGRPCHandler creates a new gRPC handler
func NewGRPCHandler(logger *zap.Logger) *GRPCHandler {
	return &GRPCHandler{
		logger: logger,
	}
}

// GetAlert returns the GetAlertResponse to the gRPC client that called it
func (g *GRPCHandler) GetAlert(ctx context.Context, r *api_v1.GetAlertRequest) (*api_v1.GetAlertResponse, error) {
	alert := make(map[string]string)
	g.logger.Info("Handling request for Get Alert..This is a test for demo purposes")
	alert["alert"] = "details"
	return &api_v1.GetAlertResponse{Alert: alert}, nil
}
