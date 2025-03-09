# ADR 0008: Use FluxCD for GitOps-Based Deployment

## Status
Accepted

## Context
To ensure **reproducible, declarative, and automated deployments**, Butler requires a GitOps-based approach. Traditional deployment methods introduce:
- **Configuration drift** due to manual changes.
- **Limited rollback capabilities** for failed deployments.
- **Complexity in multi-cluster configuration management.**

FluxCD is a **GitOps operator for Kubernetes** that provides:
- **Automated reconciliation** between Git repositories and Kubernetes state.
- **Multi-cluster and multi-tenant support** for complex environments.
- **Seamless integration with Helm and Kustomize.**

## Decision
We will use **FluxCD** as the GitOps framework for Butler because:
- **It continuously applies changes from Git**, enforcing declarative configuration.
- **Supports Helm and Kustomize**, aligning with Butler’s modular approach.
- **Improves disaster recovery** by ensuring infrastructure can be re-provisioned from Git.

## Consequences
- **Stronger deployment consistency**: Reduces configuration drift across environments.
- **Improved auditability**: Git history provides a clear record of all changes.
- **Learning curve**: Requires teams to fully adopt GitOps practices.

## Next Steps
- Deploy FluxCD in Butler’s management cluster.
- Define Git repositories and structure for cluster configurations.
- Automate deployment workflows for infrastructure and applications.

