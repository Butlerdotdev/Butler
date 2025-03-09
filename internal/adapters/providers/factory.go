package providers

import (
	"butler/internal/adapters/providers/nutanix"
	"context"
	"errors"

	"go.uber.org/zap"
)

// Supported providers
const (
	ProviderNutanix = "nutanix"
	ProviderAWS     = "aws"
	ProviderAzure   = "azure"
)

// NewProviderFactory returns a cloud provider adapter.
func NewProviderFactory(ctx context.Context, providerType string, config map[string]string, logger *zap.Logger) (ProviderInterface, error) {
	switch providerType {
	case ProviderNutanix:
		client := nutanix.NewNutanixClient(ctx, config["endpoint"], config["username"], config["password"], logger)
		return nutanix.NewNutanixAdapter(client, logger), nil
	default:
		return nil, errors.New("unsupported provider: " + providerType)
	}
}
