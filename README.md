[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![made-with-Go](https://img.shields.io/badge/Made%20with-Go-1f425f.svg)](http://golang.org)
[![Go Report Card](https://goreportcard.com/badge/github.com/Butlerdotdev/Butler)](https://goreportcard.com/report/github.com/Butlerdotdev/Butler)
[![GitHub go.mod Go version of a Go module](https://img.shields.io/github/go-mod/go-version/gomods/athens.svg)](https://github.com/gomods/athens)
[![GoDoc reference example](https://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/nanomsg.org/go/mangos/v2)

# Butler: Kubernetes as a Service Platform

## Table of Contents
- [Overview](#overview)
- [Key Features](#key-features)
- [Project Structure](#project-structure)
- [Documentation](#documentation)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
- [Contributing](#contributing)
- [License](#license)

## Overview
**Butler** is a Kubernetes-native cloud platform designed to **provision, manage, and automate** Kubernetes clusters and virtualized workloads. Butler leverages modern **cloud-native technologies** such as **Kamaji, KubeVirt, Kube-OVN, LINSTOR, Nutanix CSI, MetalLB, Traefik, and FluxCD** to provide a robust, scalable infrastructure solution.

## Key Features
- **Lightweight Control Planes** with Kamaji.
- **Unified Virtual Machine & Container Orchestration** via KubeVirt.
- **Advanced Networking** using Kube-OVN & MetalLB.
- **High-Availability Storage** with LINSTOR & Nutanix CSI.
- **Declarative, GitOps-Driven Deployment** using FluxCD & Helm.
- **Multi-Cluster & Multi-Tenant Support** for large-scale deployments.

## Project Structure
```
├── cmd/              # CLI Commands for Butler
├── internal/         # Core Logic & Adapters
│   ├── adapters/    # Integrations with Infrastructure Providers
│   ├── services/    # Business Logic for Cluster & VM Provisioning
│   ├── handlers/    # Request Handlers
├── docs/            # Documentation & ADRs
│   ├── adr/        # Architecture Decision Records
│   ├── tdd/        # Technical Design Documents
│   ├── roadmap/    # Project Roadmap
├── pkg/            # Shareble packages
└── README.md        # Project Overview
```

## Documentation
📖 **[Technical Design Document (TDD)](docs/technical-design-documents/TDD.md)**
📌 **[Project Roadmap](docs/roadmap/README.md)**
📜 **[Architecture Decision Records (ADRs)](docs/architecture-design-records/adr0000-adr-log.md)**

## Getting Started
### Prerequisites
- **Kubernetes 1.24+**
- **Helm 3+**
- **kubectl**
- **talosctl**

### Installation

## Contributing
We welcome contributions! Please review our [Contributing Guide](CONTRIBUTING.md) and open issues or pull requests.

## License
📄 **Apache License 2.0** - See [LICENSE](LICENSE) for details.

