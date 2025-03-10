package bootstrap

import (
	"butler/internal/adapters/providers"
	"butler/internal/models"
	"fmt"
	"go.uber.org/zap"
)

// Provisioner handles VM provisioning.
type Provisioner struct {
	provider providers.ProviderInterface
	logger   *zap.Logger
}

// NewProvisioner initializes a Provisioner.
func NewProvisioner(provider providers.ProviderInterface, logger *zap.Logger) *Provisioner {
	return &Provisioner{
		provider: provider,
		logger:   logger,
	}
}

// ProvisionVMs provisions the required VMs for the management cluster.
func (p *Provisioner) ProvisionVMs(config *models.BootstrapConfig) error {
	for _, node := range config.ManagementCluster.Nodes {
		for i := 1; i <= node.Count; i++ {
			vmName := fmt.Sprintf("%s-%s-%d", config.ManagementCluster.Name, node.Role, i)
			p.logger.Info("Creating VM", zap.String("name", vmName), zap.String("role", node.Role))

			vmConfig := models.VMConfig{
				Name:        vmName,
				CPU:         node.CPU,
				RAM:         node.RAM,
				Disk:        node.Disk,
				IsoUUID:     node.IsoUUID,
				SubnetUUID:  config.ManagementCluster.Nutanix.SubnetUUID,
				ClusterUUID: config.ManagementCluster.Nutanix.ClusterUUID,
			}

			_, err := p.provider.CreateVM(vmConfig)
			if err != nil {
				p.logger.Error("Failed to create VM", zap.String("vm_name", vmName), zap.Error(err))
				return err
			}
		}
	}
	return nil
}
