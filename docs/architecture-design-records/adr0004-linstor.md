# ADR 0004: Use LINSTOR for Kubernetes Storage Management

## Status
Accepted

## Context
Butler requires a **highly available, scalable, and performant storage solution** for tenant clusters. Traditional storage solutions have limitations in:
- **Multi-node data replication**.
- **Automated provisioning and dynamic storage management**.
- **Seamless Kubernetes integration**.

LINSTOR, an open-source **software-defined storage (SDS) solution**, provides:
- **Fast and efficient block storage** for Kubernetes workloads.
- **Integrated replication and failover** using DRBD technology.
- **Dynamic provisioning via CSI (Container Storage Interface).**

## Decision
We will use **LINSTOR** as the default storage backend for Butler because:
- **It supports automated volume provisioning** via CSI integration.
- **Provides HA (High Availability) storage** across multiple nodes.
- **Optimized for Kubernetes**, making it a natural fit for Butler’s architecture.

## Consequences
- **Improved storage resilience**: Automatic failover ensures data durability.
- **Better Kubernetes-native storage management**: Seamless PVC integration.
- **Operational overhead**: Butler admins must monitor LINSTOR’s replication policies and node health.

## Next Steps
- Deploy LINSTOR as the default storage backend for Butler.
- Configure LINSTOR’s CSI driver for persistent storage.
- Integrate LINSTOR with Butler’s backup and recovery policies.

