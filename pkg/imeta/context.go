package imeta

import (
	"context"
)

type mdKey struct{}

// WithContext creates a new context with md attached.
func WithContext(ctx context.Context, md MD) context.Context {
	return context.WithValue(ctx, mdKey{}, md)
}

// FromContext returns the incoming metadata in ctx if it exists.  The
// returned MD should not be modified. Writing to it may cause races.
// Modification should be made to copies of the returned MD.
func FromContext(ctx context.Context) (md MD, ok bool) {
	if ctx == nil {
		return MD{}, false
	}

	md, ok = ctx.Value(mdKey{}).(MD)
	return
}
