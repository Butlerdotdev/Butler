# ADR 0001: Use Kamaji for Lightweight Control Planes

## Status
Accepted

## Context
In Butler, we need a scalable, efficient way to manage control planes for multiple downstream Kubernetes clusters. Traditional Kubernetes deployments require running full control plane components for each cluster, which leads to:

- **High resource consumption** per cluster.
- **Operational complexity** in maintaining multiple control planes.
- **Scalability limitations** due to control plane overhead.

Kamaji is a solution that enables **lightweight Kubernetes control planes** by running them as pods inside a management cluster. This approach aligns well with Butler’s goal of managing multiple clusters efficiently.

## Decision
We will use **Kamaji** as the control plane provisioning solution for Butler. It will allow us to:

- **Reduce resource overhead** by running control planes as pods.
- **Scale efficiently** by provisioning tenant clusters with minimal control plane footprint.
- **Improve automation** through integration with Cluster API (CAPI).

## Consequences
- **Lower resource usage**: Control planes will be containerized, reducing dedicated VM/node requirements.
- **Improved manageability**: A single management cluster will orchestrate multiple lightweight control planes.
- **Dependency on Kamaji’s roadmap**: Future changes in Kamaji could impact Butler’s architecture.

## Next Steps
- Implement Kamaji in the Butler management cluster.
- Integrate Kamaji with Butler’s cluster provisioning workflows.
- Define Kamaji-specific configurations and security policies.