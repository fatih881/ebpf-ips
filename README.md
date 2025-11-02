# About The Project
This repository documents my research on an IPS (Intrusion Prevention System) powered by the combination of XDP and DDP (Dynamic Device Personalization).
### Roadmap
* **Phase 1 (software-only)**: In this phase, the functions will be coded to work all on host cpu. The main goal in this phase is developing the IPS to work without a DDP profile. In this phase I'm planning on coding other necessary implementations. (e.g., Prometheus implementation)
>During this software-only phase, the XDP program cannot rely on DDP to parse the info we specified (like custom tunnels or L7 payloads). This parsing process will be a high compute-load for a normal CPU, but with an ASIC it will be in nanoseconds. To access this data, the XDP program must perform manual, software-based parsing (pointer arithmetic), which is CPU-intensive. This slow, software-based parsing in Phase 1 is intentional. It creates the critical baseline we will benchmark against the fast, hardware-offloaded parsing of Phase 2 (DDP).
* **Phase 2 (hardware-accelerated)**: This phase is contingent on securing a NIC with XDP_NATIVE and DDP profile support. Upon acquisition, we will refactor the XDP code to use bpf_flow_dissector instead of all of the complexity we have in phase 1 on reading the HTTP payload, and we will offload the CPU load to the ASIC. With the comparative benchmark of phase 1 and 2, we will achieve the final project goal.
### Project scope and the idea
* The idea behind starting this project is firstly my self-learning progress on XDP and DDP. I also want to make comprehensive, comparative benchmarking on it. (e.g., benchmarks: since we have a lot of benchmarks for xdp/userspace, the benchmarks will contain XDP vs. XDP+ DDP)

### Goals and Non-Goals

* **What this is for:** To explore the boundaries of an IPS system with DDP+XDP. Since DDP requires more than a consumer-level NIC, DDP is a future plan and we must be content with unit and end-to-end tests.
>HTTPS needs a SuperNIC so we will be using HTTP traffic for proving the concept.
* **What this is not for:** This project is not an attempt to replace something in the industry. It's an architectural exploration.
>I am open to feedback, ideas, and discussions regarding this approach. Please take a look at contributing.md.

