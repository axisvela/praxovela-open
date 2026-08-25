# Syndovela Operational Guide (Open)

## Getting Started

### Prerequisites
- PRAXOVELA v1.9.0 or later
- HTTP client (curl, wget, HTTPie, etc.)
- Optional: Syndovela control plane configuration

### Installation

#### Via Go Module
```bash
go get github.com/axisrobo/praxovela@v1.3.0
```

#### Binary Distribution
```bash
# Download pre-built binary
# Or build from source
go build ./...
```

### Running the Facade

```go
import (
    "github.com/axisrobo/praxovela/adapters/syndovela"
    "github.com/axisrobo/praxovela/axon-core/internal/config"
    "net/http"
)

func main() {
    // Initialize facade
    descriptor := syndovela.RuntimeDescriptor{
        RuntimeID:             "praxovela",
        Implementation:        "axon-core",
        ImplementationVersion: "1.3.0",
        ProtocolVersions:      []string{"sbrp/v1"},
        Isolation:             []string{"wasm"},
        ABIs:                  []string{"wasi/preview2"},
        Platform:              "windows/amd64",
    }
    
    // Load locks from configuration
    locks := [...]syndovela.ResolutionLock{
        {LockID: "lock-1", ResolverVersion: "r1", Digest: "sha256:manifest", Selected: []syndovela.SelectedBundle{
            {BundleID: "bundle-a", Version: "1.0.0", Digest: "sha256:a"},
        }},
    }
    
    // Create facade with optional persistence
    facade := syndovela.NewRuntimeFacade(descriptor, locks, syndovela.WithPersistence("/path/to/state.json"))
    
    // Start HTTP server
    mux := http.NewServeMux()
    mux.Handle("/syndovela", facade)
    
    httpServer := &http.Server{
        Addr: ":8765",
        Handler: mux,
    }
    
    log.Printf("Syndovela facade starting on :8765")
    httpServer.ListenAndServe()
}
```

### Environment Configuration

| Variable | Description | Default |
|----------|-------------|---------|
| `AXON_SYNDOVELA_URL` | Syndovela API base URL | (empty = standalone mode) |
| `SYNDOVELA_PERSISTENCE_PATH` | Path to state JSON file | (disabled) |

### Standalone vs Wired Mode

#### Standalone Mode
```go
// No URL configured - facade runs independently
facade := syndovela.NewRuntimeFacade(descriptor, locks)
// No AXON_SYNDOVELA_URL needed
```

#### Wired Client Mode
```go
// Connect to control plane
import "github.com/axisrobo/praxovela/axon-core/cmd/axond"

// In axond main.go:
runtime.WithSyndovelaClient(syndovela.NewClient("http://127.0.0.1:8765/syndovela"))
```

## Daily Operations

### Check Capabilities
```bash
curl -X GET http://localhost:8765/syndovela/describe
```

### Apply a New Bundle
```bash
curl -X POST http://localhost:8765/syndovela/apply \
  -H "Content-Type: application/json" \
  -d '{"bundleId":"my-skill","version":"2.0.0","digest":"sha256:abc123","resolutionLockRef":"my-lock"}'
```

### Check Instance Status
```bash
curl -X GET http://localhost:8765/syndovela/report
```

### Revoke a Bundle
```bash
curl -X POST http://localhost:8765/syndovela/revoke \
  -H "Content-Type: application/json" \
  -d '{"instanceId":"inst-1"}'
```

### S7 Lifecycle Operations

#### Fetch a Bundle
```bash
curl -X POST http://localhost:8765/syndovela/fetch \
  -H "Content-Type: application/json" \
  -d '{"bundleId":"my-skill"}'
```

#### Validate a Bundle
```bash
curl -X POST http://localhost:8765/syndovela/validate \
  -H "Content-Type: application/json" \
  -d '{"instanceId":"inst-1"}'
```

#### Load a Bundle
```bash
curl -X POST http://localhost:8765/syndovela/load \
  -H "Content-Type: application/json" \
  -d '{"instanceId":"inst-1"}'
```

## Troubleshooting

### Common Issues

| Symptom | Cause | Solution |
|---------|-------|----------|
| `409 Conflict` on /apply | Digest mismatch with lock | Verify resolutionLockRef and digest match |
| `404 Not Found` on /revoke | Instance ID doesn't exist | Check /report to list active instances |
| `405 Method Not Allowed` | Wrong HTTP method | Use correct method (GET/POST as specified) |
| State stuck in DRAINING | In-flight invocations not completed | Wait for invocations to complete, or use revoke with wait flag |
| `500 Internal Server Error` | Persistence file issues | Ensure write permissions, valid JSON format |

### Debug Mode

```bash
# Enable verbose output
curl -v -X GET http://localhost:8765/syndovela/describe

# Check server logs
# (See axond logs for detailed error information)
```

## Upgrading

### v1.2.0 → v1.3.0

1. **No breaking changes** - All v1.2.0 endpoints remain compatible
2. **New endpoints**: `/fetch`, `/validate`, `/load` available
3. **Optional persistence**: Use `WithPersistence()` to enable state saving
4. **Rebuild**: `go build ./...` and restart service

### Migration Path

Existing installations can upgrade without code changes. The v1.3.0 facade
maintains full backward compatibility with v1.2.0 APIs.

## Getting Help

- **GitHub Issues**: https://github.com/axisrobo/praxovela/issues
- **Documentation**: https://github.com/axisrobo/praxovela/wiki
- **Version**: v1.3.0 (2026-08-24)

---

## License

Apache-2.0 OR MIT (per project dual-license)
See LICENSE file for details.