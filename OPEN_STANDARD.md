# The Nexus Open Standard (NOS)

Welcome to the **Nexus Open Standard**. This specification enables third-party security vendors (EDRs, SIEMs, Vulnerability Scanners) to natively integrate with the Nexus Node ecosystem via high-performance gRPC plugins.

By adopting this standard, you no longer need to deploy a proprietary agent. You can leverage the ubiquitous footprint of the Nexus Node to stream telemetry or issue remediation actions directly to endpoints.

## Core Concepts

1. **NOSEvent**: Every piece of telemetry (process execution, network connection, vulnerability scan) is packaged into a standard `NOSEvent`.
2. **Plugins**: A plugin is a binary running alongside Nexus (or deployed as a sidecar container) that speaks the NOS gRPC Protocol.
3. **Bi-Directional Action**: Plugins can not only *listen* to data but can also *command* the Agent (e.g., kill a malicious process or apply an OS patch) using the `ActionService`.

## The Interfaces

> See `pkg/nos/proto/nos_schema.proto` and `pkg/nos/proto/nos_service.proto` for the official definitions.

### Subscribing to Data (StreamService)
To receive data, a vendor plugin calls `Subscribe` on the Agent.

```go
// Example Vendor Plugin Code
stream, err := client.Subscribe(ctx, &nos.StreamRequest{
    PluginName: "Crowdstrike_Falcon_Analytics",
    FilterQuery: "event.type == 'process'", // Only stream process executions
})
for {
    event, err := stream.Recv()
    // Analyze NOSEvent in your vendor backend
}
```

### Issuing Commands (ActionService)
If a vendor detects a threat, they can command the local Nexus Agent to remediate it.

```go
// Kill a malicious process
resp, err := client.ExecuteAction(ctx, &nos.ActionRequest{
    Action:      nos.ActionType_KILL_PROCESS,
    TargetId:    "4192", // The malicious PID
    Authorization: "<vendor-signed-jwt>"
})
```

## Security & Governance
- All gRPC connections between Nexus and Plugins must be secured via **mTLS**.
- Plugins must present a strong cryptographic identity (JWT or x509 cert).
- The Nexus Node's local `pipeline.Router` (or the central Control Plane) acts as the final gatekeeper and can drop vendor commands if they violate enterprise policy.

## Getting Started
To build a Nexus Plugin:
1. Generate the Go/Rust/C++ stubs using `protoc` from the `.proto` files in `pkg/nos/proto/`.
2. Implement your processing logic.
3. Register your plugin with the Nexus Central Management Server to receive authorized access to agent data streams.
