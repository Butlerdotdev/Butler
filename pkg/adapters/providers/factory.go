package providers

import (
	"context"
	"errors"

	"github.com/butlerdotdev/butler/pkg/adapters/providers/nutanix"
	"github.com/butlerdotdev/butler/pkg/adapters/providers/proxmox"

	"go.uber.org/zap"
)

// Supported providers
const (
	ProviderNutanix = "nutanix"
	ProviderAWS     = "aws"
	ProviderAzure   = "azure"
	ProviderProxmox = "proxmox"
)

// NewProviderFactory returns a cloud provider adapter.
func NewProviderFactory(ctx context.Context, providerType string, config map[string]string, logger *zap.Logger) (ProviderInterface, error) {
	switch providerType {
	case ProviderNutanix:
		client := nutanix.NewNutanixClient(ctx, config["endpoint"], config["username"], config["password"], logger)
		return nutanix.NewNutanixAdapter(client, logger), nil
	case ProviderProxmox:
		client := proxmox.NewProxmoxClient(ctx, config["endpoint"], config["username"], config["password"], config["nodes"], logger)
		return proxmox.NewProxmoxAdapter(client, logger), nil
	default:
		return nil, errors.New("unsupported provider: " + providerType)
	}
}
