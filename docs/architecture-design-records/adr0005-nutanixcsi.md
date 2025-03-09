# ADR 0005: Use Nutanix CSI for Storage Integration

## Status
Accepted

## Context
Butler needs **seamless storage integration** with Nutanix AHV to support workloads running on Nutanix infrastructure. Kubernetes requires a **CSI (Container Storage Interface) driver** to dynamically provision persistent storage.

The **Nutanix CSI driver** provides:
- **Automated volume provisioning** from Nutanix storage pools.
- **Snapshot and backup support** for Kubernetes workloads.
- **Integration with Nutanix Prism** for unified storage management.

## Decision
We will use **Nutanix CSI** as a storage backend in Butler where Nutanix AHV is the underlying infrastructure because:
- **It natively integrates with Nutanix’s hyperconverged storage.**
- **Supports dynamic PVC provisioning** for Kubernetes applications.
- **Provides high availability and snapshot support.**

## Consequences
- **Seamless storage experience**: Workloads on Nutanix clusters benefit from native storage integration.
- **Tight Nutanix dependency**: Workloads relying on Nutanix CSI must run on Nutanix infrastructure.
- **Operational complexity**: Requires ongoing monitoring of Nutanix storage pools.

## Next Steps
- Deploy the Nutanix CSI driver in Butler’s Kubernetes environment.
- Configure storage classes to support dynamic provisioning.
- Define Nutanix storage policies for workload resilience.

