# Contributing to Nexus Node

We love your input! We want to make contributing to this project as easy and transparent as possible.

## Development Workflow
1.  **Fork** the repo on GitHub.
2.  **Clone** the project to your own machine.
3.  **Branch** based on the feature you are adding (e.g., `feature/linux-ebpf`).
4.  **Commit** changes to your own branch.
5.  **Push** your work back to your fork.
6.  Submit a **Pull Request** so that we can review your changes.

## Module Architecture
New capabilities should be implemented as **Modules** in `/internal/modules`.
-   Must implement a `Start()` / `Stop()` lifecycle if long-running.
-   Should respect the `DryRun` flag where applicable.
-   Must output data in **OCSF** format if generating telemetry.

## Style Guide
-   Use `go fmt` before committing.
-   Comments should be full sentences.
-   Code must pass `go vet`.

## Reporting Issues
We use GitHub issues to track public bugs. Report a bug by opening a new issue; it's that easy!
