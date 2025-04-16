package providers

import "github.com/butlerdotdev/butler/pkg/models"

// ProviderInterface defines required cloud provider operations.
type ProviderInterface interface {
	CreateVM(vm models.VMConfig) (string, error)
	DeleteVM(vmID string) error
	GetVMStatus(vmName string) (models.VMStatus, error)
}
