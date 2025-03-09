# ADR 0003: Use Kube-OVN for Kubernetes Networking

## Status
Accepted

## Context
Networking in Butler must be **scalable, secure, and multi-tenant aware** to support various workloads across multiple clusters. Kubernetes’ default networking implementations lack:
- **Built-in multi-tenancy and network isolation**.
- **Advanced traffic control and observability**.
- **Seamless integration with hybrid cloud infrastructures**.

Kube-OVN is an **open-source Kubernetes networking solution** that extends the Open vSwitch (OVS) framework to provide:
- **Native Kubernetes CNI support** with flexible IP management.
- **Multi-tenancy and security policies** for strict workload isolation.
- **Enhanced network performance** via OVS-based acceleration.

## Decision
We will use **Kube-OVN** as the primary CNI (Container Network Interface) for Butler because:
- **Provides built-in security policies** for isolating tenant workloads.
- **Integrates with Kubernetes natively**, simplifying network operations.

## Consequences
- **Improved multi-tenancy**: Tenant clusters will have stronger network isolation.
- **Better observability**: Kube-OVN offers enhanced monitoring and traffic control.
- **Requires operational expertise**: Butler admins must manage OVS and advanced networking policies.

## Next Steps
- Deploy Kube-OVN as Butler’s default CNI.
- Configure network isolation policies for tenant clusters.
- Integrate Kube-OVN with Butler’s monitoring and observability stack.