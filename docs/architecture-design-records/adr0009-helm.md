# ADR 0009: Use Helm for Kubernetes Package Management

## Status
Accepted

## Context
Managing Kubernetes manifests across multiple clusters requires a **scalable and maintainable** approach. Traditional methods of applying raw YAML files introduce challenges such as:
- **Manual complexity** in updating deployments.
- **Inconsistent configurations** across environments.
- **Limited versioning and rollback capabilities.**

Helm is a **package manager for Kubernetes** that provides:
- **Templated manifests** to enable configuration reuse.
- **Version-controlled releases** for managing application lifecycles.
- **Integrated dependency management** for complex deployments.

## Decision
We will use **Helm** as the package manager for Butler because:
- **It simplifies application deployment** with reusable Helm charts.
- **Supports version-controlled rollouts**, reducing risk in production.
- **Integrates seamlessly with FluxCD**, ensuring GitOps-driven deployment consistency.

## Consequences
- **Standardized deployment process**: Helm provides a structured approach to managing Kubernetes resources.
- **Improved rollback and upgrade mechanisms**: Helm’s release tracking enables safer updates.
- **Additional learning curve**: Teams must understand Helm chart structure and Helm-specific templating.

## Next Steps
- Define Helm chart structure for Butler’s core services.
- Integrate Helm with FluxCD for automated GitOps deployments.
- Establish Helm release management policies for versioning and rollback strategies.