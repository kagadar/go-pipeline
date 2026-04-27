package seq

import (
	"iter"
	"slices"
)

type iterSeq[E any] = func(func(E) bool)

// All returns whether all elements of the provided [iter.Seq] satisfy the provided function.
func All[I ~iterSeq[E], E any](i I, f func(E) bool) bool {
	for e := range i {
		if !f(e) {
			return false
		}
	}
	return true
}

// Any returns whether any element of the provided [iter.Seq] satisfies the provided function.
func Any[I ~iterSeq[E], E any](i I, f func(E) bool) bool {
	for e := range i {
		if f(e) {
			return true
		}
	}
	return false
}

// Chunk returns a [iter.Seq] of slices of elements of the provided [iter.Seq].
// All but the last slice will contain exactly size elements.
// The last slice will be clipped to exactly the number of elements in the slice.
// Chunk panics if n is less than 1.
func Chunk[I ~iterSeq[E], E any](i I, size int) iter.Seq[[]E] {
	if size < 1 {
		panic("size cannot be less than 1")
	}
	return func(yield func([]E) bool) {
		chunk := make([]E, 0, size)
		for e := range i {
			chunk = append(chunk, e)
			if len(chunk) == size {
				if !yield(chunk) {
					return
				}
				chunk = make([]E, 0, size)
			}
		}
		if len(chunk) > 0 {
			yield(slices.Clip(chunk))
		}
	}
}

// Concat returns a [iter.Seq] that yields all elements from the provided [iter.Seq]s in order.
func Concat[I ~iterSeq[E], E any](i ...I) iter.Seq[E] {
	return func(yield func(E) bool) {
		for _, s := range i {
			for e := range s {
				if !yield(e) {
					return
				}
			}
		}
	}
}

// Contains returns whether the provided [iter.Seq] contains the specified element.
func Contains[I ~iterSeq[E], E comparable](i I, target E) bool {
	return Any(i, func(e E) bool { return e == target })
}

// Count returns the number of elements in the provided [iter.Seq].
func Count[I ~iterSeq[E], E any](i I) (o int) {
	for range i {
		o++
	}
	return
}

// Dedupe returns a [iter.Seq] over the distinct elements of the provided [iter.Seq], in order of first occurrence.
func Dedupe[I ~iterSeq[E], E comparable](i I) iter.Seq[E] {
	return func(yield func(E) bool) {
		m := map[E]struct{}{}
		for e := range i {
			if _, ok := m[e]; !ok {
				m[e] = struct{}{}
				if !yield(e) {
					return
				}
			}
		}
	}
}

// Filter returns a [iter.Seq] over all of the elements for which the provided function returned true.
func Filter[I ~iterSeq[E], E any](i I, f func(E) bool) iter.Seq[E] {
	return func(yield func(E) bool) {
		for e := range i {
			if f(e) && !yield(e) {
				return
			}
		}
	}
}

// Find returns the first element of the provided [iter.Seq] that satisfies the provided function, and true.
// If no elements satisfy the function, it returns the zero value and false.
func Find[I ~iterSeq[E], E any](i I, f func(E) bool) (e E, ok bool) {
	for e := range i {
		if f(e) {
			return e, true
		}
	}
	return
}

// Flatten returns a [iter.Seq] over all of the elements of the slices in the provided [iter.Seq].
func Flatten[I ~iterSeq[S], S ~[]E, E any](i I) iter.Seq[E] {
	return func(yield func(E) bool) {
		for s := range i {
			for _, e := range s {
				if !yield(e) {
					return
				}
			}
		}
	}
}

// FlattenSeq returns a [iter.Seq] over all of the elements of the [iter.Seq]s in the provided [iter.Seq].
func FlattenSeq[I ~iterSeq[S], S ~iterSeq[E], E any](i I) iter.Seq[E] {
	return func(yield func(E) bool) {
		for s := range i {
			for e := range s {
				if !yield(e) {
					return
				}
			}
		}
	}
}

// GroupBy returns a map of the elements from the provided [iter.Seq] grouped by the value returned by running the key func on them.
func GroupBy[I ~iterSeq[E], K comparable, E any](i I, f func(E) K) map[K][]E {
	o := map[K][]E{}
	for e := range i {
		k := f(e)
		o[k] = append(o[k], e)
	}
	return o
}

// Last returns the last element of the provided [iter.Seq].
func Last[I ~iterSeq[E], E any](i I) (o E) {
	for o = range i {
	}
	return
}

// Partition returns two new slices of the elements of the provided [iter.Seq], splitting them based on whether the provided func returned true.
func Partition[I ~iterSeq[E], E any](i I, f func(E) bool) (pass, fail []E) {
	for e := range i {
		if f(e) {
			pass = append(pass, e)
		} else {
			fail = append(fail, e)
		}
	}
	return
}

// Reduce runs the provided function for each element of the provided [iter.Seq], accumulating the result in `O`.
func Reduce[I ~iterSeq[E], O, E any](i I, f func(O, E) O) (o O) {
	for e := range i {
		o = f(o, e)
	}
	return
}

// Skip returns a [iter.Seq] over the elements of the provided [iter.Seq] after having skipped n entries.
func Skip[I ~iterSeq[E], E any](i I, n int) iter.Seq[E] {
	return func(yield func(E) bool) {
		var count int
		for e := range i {
			if count < n {
				count++
				continue
			}
			if !yield(e) {
				return
			}
		}
	}
}

// SkipWhile returns a [iter.Seq] over the elements of the provided [iter.Seq], skipping elements from the start until provided function returns true for the first time.
func SkipWhile[I ~iterSeq[E], E any](i I, f func(E) bool) iter.Seq[E] {
	return func(yield func(E) bool) {
		skipping := true
		for e := range i {
			if skipping {
				if f(e) {
					continue
				}
				skipping = false
			}
			if !yield(e) {
				return
			}
		}
	}
}

// Take returns a [iter.Seq] over at most n elements from the provided [iter.Seq].
func Take[I ~iterSeq[E], E any](i I, n int) iter.Seq[E] {
	return func(yield func(E) bool) {
		var count int
		for e := range i {
			if count >= n || !yield(e) {
				return
			}
			count++
		}
	}
}

// TakeWhile returns a [iter.Seq] over the elements of the provided [iter.Seq] until the provided function returns false.
func TakeWhile[I ~iterSeq[E], E any](i I, f func(E) bool) iter.Seq[E] {
	return func(yield func(E) bool) {
		for e := range i {
			if !f(e) || !yield(e) {
				return
			}
		}
	}
}

// ToSeq2 uses the provided function to transform the elements of the provided [iter.Seq] into a [iter.Seq2].
func ToSeq2[I ~iterSeq[E], K, V, E any](i I, f func(E) (K, V)) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for e := range i {
			if !yield(f(e)) {
				return
			}
		}
	}
}

// Transform uses the provided function to transform the elements of the provided [iter.Seq] into a new [iter.Seq].
func Transform[I ~iterSeq[E1], E1, E2 any](i I, f func(E1) E2) iter.Seq[E2] {
	return func(yield func(E2) bool) {
		for e := range i {
			if !yield(f(e)) {
				return
			}
		}
	}
}

// TransformErr uses the provided function to transform the elements of the provided [iter.Seq] into a new [iter.Seq].
// When the first non-nil error is encountered this function will cease iteration and set errOut to the error value.
func TransformErr[I ~iterSeq[E1], E1, E2 any](i I, f func(E1) (E2, error), errOut *error) iter.Seq[E2] {
	return func(yield func(E2) bool) {
		for ei := range i {
			eo, err := f(ei)
			if err != nil {
				*errOut = err
				return
			}
			if !yield(eo) {
				return
			}
		}
	}
}
