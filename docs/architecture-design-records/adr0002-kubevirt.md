# ADR 0002: Use KubeVirt for Virtual Machine Management

## Status
Accepted

## Context
Butler aims to support **both containerized and virtualized workloads**. Since Kubernetes natively orchestrates containers, we need an efficient way to manage virtual machines (VMs) alongside Kubernetes pods.

KubeVirt is an open-source project that enables running **VMs as native Kubernetes resources**, providing:
- A unified control plane for **both VMs and containers**.
- **Seamless VM lifecycle management** inside Kubernetes.
- **Flexibility** to support legacy applications that require virtualization.

## Decision
We will use **KubeVirt** to manage VMs within Butler. This decision is based on:

- **Integration with Kubernetes**: KubeVirt extends Kubernetes with a CRD-based VM abstraction.
- **Operational efficiency**: VMs can be managed using standard Kubernetes workflows.
- **Multi-tenancy support**: KubeVirt aligns with Butler’s need to provision VMs dynamically.

## Consequences
- **Increased flexibility**: Allows Butler to orchestrate **VMs and containers** side by side.
- **Improved standardization**: Unifies VM lifecycle management under Kubernetes-native APIs.
- **Potential performance considerations**: Running VMs inside Kubernetes introduces additional overhead.

## Next Steps
- Deploy KubeVirt in the Butler management cluster.
- Define VM provisioning workflows within Butler.
- Integrate KubeVirt with Butler’s cluster API-based provisioning system.