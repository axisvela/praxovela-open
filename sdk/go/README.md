# PRAXOVELA Go SDK

Module: `github.com/axisvela/praxovela-open/sdk/go/v3`

```go
import sdk "github.com/axisvela/praxovela-open/sdk/go/v3/praxovela"

c, err := sdk.New(sdk.Config{BaseURL: "http://127.0.0.1:1926", Token: "..."})
if err != nil { /* ... */ }

sess, err := c.CreateSession(ctx, "workspace-1")
run, err := c.StartRun(ctx, sess.ID, "review this change")
status, err := c.GetRun(ctx, run.ID)
err = c.CancelRun(ctx, run.ID)
recovery, err := c.RecoverRun(ctx, run.ID)
```

The SDK wraps the PRAXOVELA/AXON Core HTTP API; callers should not hand-roll
HTTP for these operations.
