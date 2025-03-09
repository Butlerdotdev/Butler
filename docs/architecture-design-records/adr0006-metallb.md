# ADR 0006: Use MetalLB for Load Balancing in Bare Metal Environments

## Status
Accepted

## Context
Kubernetes does not provide a built-in **LoadBalancer** for bare metal environments. In Butler, we need a way to assign **external IPs to services** in a cluster running on Nutanix AHV or other bare metal platforms.

MetalLB is an open-source **network load balancer** that provides:
- **External IP allocation** for Kubernetes services in a bare metal environment.
- **BGP (Border Gateway Protocol) and Layer 2 modes** for flexible network integration.
- **Seamless integration with Kubernetes services.**

## Decision
We will use **MetalLB** as Butler’s default **load balancer** in bare metal deployments because:
- **It enables LoadBalancer services** in environments without cloud provider support.
- **Supports BGP and L2 modes**, making it adaptable to different network setups.
- **Allows high availability** by distributing traffic across multiple nodes.

## Consequences
- **Enables external access to services** on bare metal Kubernetes clusters.
- **Network setup complexity**: Requires careful configuration of BGP or L2 mode.
- **Dependency on network infrastructure**: Requires cooperation with underlying networking hardware.

## Next Steps
- Deploy MetalLB in Butler’s management cluster.
- Configure MetalLB to allocate external IPs dynamically.
- Test integration with Butler’s ingress and service mesh components.

