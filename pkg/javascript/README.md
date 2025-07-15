# Javascript

Provides a simple abstraction ontop of [sobek](https://github.com/grafana/sobek) for loading user scripts and injecting global dbscript functions for interacting with events.

## runtime.go

Providers runtime functions for interacting with events.

- `dbscript.ctx.ok`
- `dbscript.ctx.drop`
- `dbscript.ctx.error`

## javascript.go

Handle to a user provided script. Initializes a VM instance with runtime functions, compiles user script and executes on new events.

Each call to `Execute` creates a new execution context with an event containing metadata and payload. Functions injected from `runtime.go` allow interacing with the execution context - modify, drop or passthrough to next step in the pipeline.
