# ADR 0007: Use Traefik as the Ingress Controller

## Status
Accepted

## Context
Kubernetes workloads require **ingress routing** to expose services to external users. Butler needs an ingress controller that supports:
- **Dynamic configuration updates** without reloading.
- **Flexible routing** with support for custom middleware.
- **Integration with MetalLB for bare metal environments.**

Traefik is a **lightweight, modern ingress controller** that provides:
- **Automatic discovery of services** and routing.
- **Built-in support for Let's Encrypt** and certificate management.
- **Seamless integration with Kubernetes and service meshes.**

## Decision
We will use **Traefik** as Butler’s default **ingress controller** because:
- **It dynamically updates routes** without requiring restarts.
- **Supports advanced traffic control** via middleware plugins.
- **Integrates with MetalLB** to provide external access.

## Consequences
- **Improved ingress management**: Services can be exposed efficiently.
- **Increased flexibility**: Traefik’s middleware allows fine-grained traffic control.
- **Learning curve**: Requires administrators to understand Traefik’s configuration model.

## Next Steps
- Deploy Traefik as Butler’s ingress controller.
- Configure Traefik for certificate management and secure ingress.
- Integrate Traefik with MetalLB and Butler’s multi-cluster networking model.

