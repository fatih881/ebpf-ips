# Distributed API-first Firewall for every edge
[![CI](https://github.com/fatih881/ebpf-fw/actions/workflows/main.yml/badge.svg)](https://github.com/fatih881/ebpf-fw/actions/workflows/main.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/fatih881/ebpf-fw)](https://goreportcard.com/report/github.com/fatih881/ebpf-fw)
[![codecov](https://codecov.io/gh/fatih881/ebpf-fw/graph/badge.svg?token=IVF4HTGMWB)](https://codecov.io/gh/fatih881/ebpf-fw)
> **Status:** Not Under Development,No Releases,Archived 
## This repository is archived as the current architecture has reached its limits in deterministic state management . The technical insights and performance bottlenecks identified here serve as the foundational research for future projects & architectures. 
### Technical Debts 
 - The state retrieval channel copies the entire map for each response. Future designs will implement another channel for availability-check mechanisms to mitigate memory overhead and latency.
 - The `attach_stress_test.go` logic suffers from a lack of atomic synchronization between user-space slice ordering and kernel link indices. Within current implementation, removing a specific range of interfaces does not guarantee consistency with the kernel state, leading to "stale handle" errors during the test.
 - Build artifacts created via Packer contain hardcoded credentials within Docker Compose files. This exposes sensitive data in the image filesystem and layer history.Also, the images can be logged in with fedora/packer.
> The issues listed above were slated for resolution; however, following the decision to archive this iteration, they remain unpatched. These vulnerabilities and logic flaws will be taken into consideration in Successor projects.

This repository hosts the source code for a low-latency, API-first distributed firewall solution designed for the linux kernel. By leveraging eBPF and XDP, this project aims to provide line-rate packet processing with minimal userspace intervention.

The goal is to move from CLI-based firewalls to API-first distributed firewalls for every edge we have.

## Core Principles

The development is following three fixed principles:

### API-First & Scalability
The philosophy behind this system is **What if this program is running on 1 million nodes?**,and the scalability both mean using the service in high scale and development scale.
* **Programmability:** The control plane is exposed entirely via gRPC APIs (Planned Architecture) and logs/messages are using prometheus/zap logger.
* **Orchestration Ready:** Designed to be integrated into larger orchestration systems (Kubernetes, Custom Control Planes) rather than manually configured by human operators.Since it's pure gRPC,users can use any data source as long as they use gRPC.

### Autonomous Operations
The whole system is designed to require minimal manual intervention for operations.Capabilities like ``XDP_Offload,XDP_DRIVER`` are managed via service and do not require any configuration.
* **Auto-Discovery:** The service automatically detects available network interfaces for network traffic.
* **Kernel-driven event track:** The service listens to RTM_NEWLINK messages from the kernel and attaches the XDP program before traffic is allowed to pass. (See [kernelsubscription.go](core/netlink/kernelsubscription.go))
### Roadmap & TODO

 - [x] Automated attaching process (Scanning all interfaces on startup,then listening for RTM_NETlINK)
   - [x] Implemention(See issue #11)
   - [x] unit tests,integration tests
 - [ ] Implementing e2e tests for popular kernel versions in production with CI 
 - [ ] gRPC API integration for firewall rule management
 - [ ] Dry-Run Mode (To prevent losing access for management in first run,the program is going to use a Dry-Run mode.This mode will log actions without affecting the real traffic.(e.g., *"Packet would be dropped"*))
> Note that repository name will be changed.
