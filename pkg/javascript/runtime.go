package javascript

import "log/slog"

// TODO: figure out a better setup to only ever require injecting events into
// the runtimeCtx without the need of recreating the context every time.
// Something that would also allow retrieving the event from the runtimeCtx handle
// for further processing in other handlers or sending to the configured "sink"

type runtimeCtx struct{}

type Runtime struct {
	Context runtimeCtx `json:"ctx"`
}

// TODO: properly type event when figured out
func (f *runtimeCtx) Ok(event any) {
	slog.Info("Ok called", event)
}

// TODO: properly type event when figured out
func (f *runtimeCtx) Drop(reason string, event any) {
	slog.Info("Drop called", reason, event)
}

// TODO: properly type event when figured out
func (f *runtimeCtx) Error(event any) {
	slog.Info("Error called", event)
}
