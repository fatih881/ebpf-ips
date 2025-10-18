# How to contribute?
### This project focuses on providing robust, low-latency resource control for network control planes on Linux. We value reliability, performance, and clear, minimal code.
## 1. Rationale and Scope

We maintain an intentionally narrow scope: Linux Kernel API's, Go, and eBPF.
- In Scope: Linux syscalls, kernel map manipulation, Go concurrency, performance optimization on bare metal.
- Out of Scope: Windows/macOS compatibility, excessive logging or tracing overhead, adding dependencies that introduce significant runtime cost.

If your contribution does not directly enhance the stability or performance of the Control Plane on Linux, it will likely be declined.

## Using host resources efficiently
Our philosophy is that the observability agent must never be the source of a resource problem. We acquire logs to troubleshoot, not to cause overhead.
- Resource Overhead: Contributions introducing significant or unwarranted CPU, Memory, or I/O load will be rejected. The performance gain must always outweigh the resource cost.
- Avoid Polling: Do not use blocking loops or excessive polling to detect status changes (e.g., continually looping to check a flag). Prefer event-driven or channel-based Go primitives to ensure efficiency.

## About Workflow and Code Quality

We demand clarity and provable quality. Submission is not complete without verification.

- Testing is Mandatory: Any code introducing new logic must be accompanied by comprehensive Unit Tests. Code without adequate test coverage will be declined.
- Maintainability: Clarity is non-negotiable. Use descriptive branch names (feat/..., fix/...), maintain concise and clear function/variable names, and ensure all error logs are unambiguous and actionable. Code that requires excessive effort to understand will not be merged.
- Issue Tracking: All functional changes are tracked via Issues. Any work that bypasses the standard workflow (Code First) must still be formally documented: An Issue is required to be opened (retrospectively) to archive the work, even if the feature is already complete.
