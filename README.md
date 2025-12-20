# Nexus Node 🕷️

> **The "Super Agent" for the modern enterprise.**
> A modular, high-performance security agent replacing Splunk Forwarders, Tanium, Qualys, and generic Patch Managers with a single open-architecture binary.

![Build Status](https://img.shields.io/badge/build-passing-brightgreen)
![Go Version](https://img.shields.io/badge/go-1.21%2B-blue)
![Platform](https://img.shields.io/badge/platform-windows%20%7C%20linux%20%7C%20macos-lightgrey)

## 🚀 Overview

Nexus Node is designed to break the vendor lock-in of proprietary agents. It embraces **Open Standards** (OCSF, eBPF) and **Best-in-Class Open Source** tools (Osquery, Trivy, Winget) to deliver a Unified Endpoint Management & Security platform.

### Core Capabilities
-   **👁️ Watcher**: Asset Inventory & config monitoring via embedded [Osquery](https://osquery.io).
-   **🏹 Hunter**: Vulnerability scanning via embedded [Trivy](https://trivy.dev) and real-time eBPF telemetry.
-   **🛠️ Doer**: Safe, policy-driven patching using native OS package managers (`winget`, `apt`, `softwareupdate`).
-   **❤️ Healer**: Auto-remediation engine that fixes critical vulnerabilities automatically based on policy.

## 🏗️ Architecture

Nexus Node follows a "Chassis" architecture. The core binary is a lightweight orchestrator that manages specialized sub-modules.

```mermaid
graph TD
    A[Nexus Chassis] -->|Manage| B(Osquery Module)
    A -->|Manage| C(Trivy Module)
    A -->|Check| D(Safety System)
    D -->|Gate| E(Patcher Module)
    A -->|Collect| F(Telemetry/eBPF)
    A -->|Export| G{Splunk HEC / OCSF}
    C -->|Findings| H[Remediation Engine]
    H -->|Trigger| E
```

## 📥 Installation

**Nexus Node** is distributed as a single static binary. You do **not** need to build it yourself.

### Quick Install (Production)

**Windows (PowerShell)**:
```powershell
Set-ExecutionPolicy Bypass -Scope Process -Force; [System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor 3072; iex ((New-Object System.Net.WebClient).DownloadString('https://raw.githubusercontent.com/YOUR_USERNAME/nexus-node/main/scripts/install.ps1'))
```

**Linux / macOS**:
```bash
curl -sfL https://raw.githubusercontent.com/YOUR_USERNAME/nexus-node/main/scripts/install.sh | sudo bash
```

## ⚡ Usage

1.  **Configure**: the installer places a default `config.yaml` in the install directory (e.g., `C:\Program Files\NexusNode`).
    ```yaml
    agent:
      name: "node-01"
    exporter:
      splunk:
        enabled: true
    remediation:
      enabled: true
    ```
2.  **Run**:
    ```bash
    # Windows
    & "C:\Program Files\NexusNode\nexus.exe"
    
    # Linux
    /opt/nexus-node/nexus
    ```
3.  **Logs**: Check functionality.
    -   Startup: `Agent Chassis initialized...`
    -   Dry Run Patching: `Patcher: ListUpdates completed...`
    -   Remediation: `[REMEDIATION] Attempting to fix...`
    
## 🤝 Contributing

We welcome contributions! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for details on how to add new modules (e.g., eBPF for Linux).

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
