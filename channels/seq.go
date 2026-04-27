package channels

import (
	"context"
	"iter"
)

// CollectSeq collects values from the provided channel until it is closed, returning those values as a [iter.Seq].
// The error will be nil unless the context is closed before the provided channel is closed, in which case it will return all of the results collected so far alongside the context's error.
func CollectSeq[I ~chan E, E any](ctx context.Context, i I, errOut *error) iter.Seq[E] {
	return func(yield func(E) bool) {
		for {
			e, ok, err := Await(ctx, i)
			if err != nil {
				*errOut = err
				return
			}
			if !ok || !yield(e) {
				return
			}
		}
	}
}
