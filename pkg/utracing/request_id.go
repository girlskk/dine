package utracing

import "context"

// contextKeyRequestID is the unique key in the context where the request ID is
// stored.
type contextKeyRequestID struct{}

// RequestIDFromContext pulls the request ID from the context, if one was set.
// If one was not set, it returns the empty string.
func RequestIDFromContext(ctx context.Context) string {
	if v, ok := ctx.Value(contextKeyRequestID{}).(string); ok {
		return v
	}
	return ""
}

// NewRequestIDContext sets the request ID on the provided context, returning a new
// context.
func NewRequestIDContext(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, contextKeyRequestID{}, id)
}
