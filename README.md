# Distributed API-first Firewall for every edge
[![CI](https://github.com/fatih881/ebpf-ips/actions/workflows/main.yml/badge.svg)](https://github.com/fatih881/ebpf-ips/actions/workflows/main.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/fatih881/ebpf-ips)](https://goreportcard.com/report/github.com/fatih881/ebpf-ips)
> **Status:** Under Development,No Releases

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