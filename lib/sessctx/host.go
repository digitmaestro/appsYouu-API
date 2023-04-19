package sessctx

import (
	"context"
)

const HostContextKey = "host"

// SetHost returns a context that includes the user value.
func SetHost(ctx context.Context, host string) context.Context {
	return context.WithValue(ctx, HostContextKey, host)
}

// Host returns the host from the context.
func Host(ctx context.Context) string {
	return ctx.Value(HostContextKey).(string)
}
