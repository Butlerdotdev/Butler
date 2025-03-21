package mappers

import "butler/internal/models"

func NutanixToMap(cfg models.NutanixConfig) map[string]string {
	return map[string]string{
		"endpoint":    cfg.Endpoint,
		"username":    cfg.Username,
		"password":    cfg.Password,
		"clusterUUID": cfg.ClusterUUID,
		"subnetUUID":  cfg.SubnetUUID,
	}
}
