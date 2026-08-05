# Sandbox Provider Interface

Stable contract between L3 Janus Gateway (Tool) and L2 Vulcan Forge (Sandbox).

## Principle

Sandbox has no knowledge of user intent. It only executes constrained specs.
High-risk tools must execute inside a sandbox.

## Go Interface

```go
type SandboxProvider interface {
    Create(ctx context.Context, spec SandboxSpec) (SandboxSession, error)
    Exec(ctx context.Context, sessionID string, command ExecSpec) (<-chan ExecEvent, error)
    CopyOut(ctx context.Context, sessionID string, paths []string) ([]Artifact, error)
    Destroy(ctx context.Context, sessionID string) error
}

type SandboxSession struct {
    ID        string
    Status    string // "running", "exited", "timed_out"
    CreatedAt time.Time
}

type SandboxSpec struct {
    Provider string            `json:"provider"`
    Image    string            `json:"image,omitempty"`
    Mounts   []MountSpec       `json:"mounts"`
    Env      map[string]string `json:"env"`
    Network  string            `json:"network"` // "disabled", "allowlisted", "full"
    Limits   ResourceLimits    `json:"limits"`
}

type MountSpec struct {
    Host  string `json:"host"`
    Guest string `json:"guest"`
    Mode  string `json:"mode"` // "ro" or "rw"
}

type ResourceLimits struct {
    MemoryMB       int `json:"memory_mb"`
    TimeoutSeconds int `json:"timeout_seconds"`
    MaxCPUPercent  int `json:"max_cpu_percent"`
}

type ExecSpec struct {
    Command string   `json:"command"`
    Args    []string `json:"args,omitempty"`
    Cwd     string   `json:"cwd"`
    TTL     int      `json:"ttl"` // seconds
}

type ExecEvent struct {
    Type    string `json:"type"` // "stdout", "stderr", "exit"
    Data    []byte `json:"data"`
    Code    int    `json:"code,omitempty"`
}

type Artifact struct {
    Path string `json:"path"`
    Size int64  `json:"size"`
    Hash string `json:"hash"`
}
```

## Default Sandbox Policy

```yaml
sandbox_default:
  network: disabled
  filesystem: scoped
  credentials: blocked
  env: sanitized
  timeout_seconds: 120
  max_memory_mb: 2048
  audit: enabled
```

## Supported Providers

| Provider | Platform | Use Case |
|----------|----------|----------|
| Docker | All (with Docker) | Container isolation, build, test |
| Podman | Linux/macOS | Rootless container alternative |
| WASM | All | Deterministic transforms, plugins |
| OS Native | Platform-specific | No-container environments |

## Implementation Status

The `Create/Exec/CopyOut/Destroy` lifecycle defined above is the **target stable protocol**. Current Go sandbox providers implement a simpler discovery interface:

- `Name() string` — provider identity for registry lookup
- `Available() bool` — runtime capability check (e.g., Docker installed)

Concrete execution uses package-specific `Execute` methods on each provider rather than the full `SandboxProvider` interface. The target protocol will be implemented as providers mature toward full lifecycle management.
