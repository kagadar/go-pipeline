package channels

import (
	"context"
	"slices"
)

// Await blocks until the provided channel returns a value or is closed.
// The error will be nil unless the context is closed before the provided channel yields, in which case it will return immediately with the context's error.
func Await[I ~chan E, E any](ctx context.Context, i I) (t E, ok bool, err error) {
	select {
	case <-ctx.Done():
		return t, false, ctx.Err()
	case t, ok = <-i:
	}
	return
}

// Collect collects values from the provided channel until it is closed, returning those values as a new slice.
// The error will be nil unless the context is closed before the provided channel is closed, in which case it will return all of the results collected so far alongside the context's error.
func Collect[O []E, I ~chan E, E any](ctx context.Context, i I) (_ O, err error) {
	return slices.Collect(CollectSeq(ctx, i, &err)), err
}
