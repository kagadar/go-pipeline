package seq

import (
	"iter"
	"maps"
	"slices"

	"github.com/kagadar/go-pipeline/must"
)

type iterSeq2[K, V any] = func(func(K, V) bool)

// All2 returns whether all pairs of the provided [iter.Seq2] satisfy the provided function.
func All2[I ~iterSeq2[K, V], K, V any](i I, f func(K, V) bool) bool {
	for k, v := range i {
		if !f(k, v) {
			return false
		}
	}
	return true
}

// Any2 returns whether any pair of the provided [iter.Seq2] satisfies the provided function.
func Any2[I ~iterSeq2[K, V], K, V any](i I, f func(K, V) bool) bool {
	for k, v := range i {
		if f(k, v) {
			return true
		}
	}
	return false
}

// AppendErr appends the elements of the provided [iter.Seq2] to the provided slice.
// When the first non-nil error is encountered a nil slice and the error will be returned.
// The provided slice may have been mutated even if an error occurred.
func AppendErr[I ~iterSeq2[E, error], O ~[]E, E any](o O, i I) (_ O, err error) {
	return must.ZeroErr(slices.AppendSeq(o, CatchErr(i, &err)), err)
}

// CatchErr returns a new [iter.Seq] of the elements in the provided [iter.Seq2].
// When the first non-nil error is encountered the [iter.Seq] will terminate and errOut will contain the error.
func CatchErr[I ~iterSeq2[E, error], E any](i I, errOut *error) iter.Seq[E] {
	return func(yield func(E) bool) {
		for e, err := range i {
			if err != nil {
				*errOut = err
				return
			}
			if !yield(e) {
				return
			}
		}
	}
}

// CollateMap collates the elements of the provided [iter.Seq2] into a new map preallocated to the provided length.
func CollateMap[O ~map[K][]V, I ~iterSeq2[K, V], K comparable, V any](length int, i I) O {
	o := make(O, length)
	for k, v := range i {
		o[k] = append(o[k], v)
	}
	return o
}

// CollectErr collects the elements of the provided [iter.Seq2] into a new slice and returns it.
// When the first non-nil error is encountered a nil slice and the error will be returned.
func CollectErr[I ~iterSeq2[E, error], E any](i I) (_ []E, err error) {
	return must.ZeroErr(slices.Collect(CatchErr(i, &err)), err)
}

// CollectMap collects the elements of the provided [iter.Seq2] into a new map preallocated to the provided length.
func CollectMap[O ~map[K]V, I ~iterSeq2[K, V], K comparable, V any](length int, i I) O {
	return InsertMap(make(O, length), i)
}

// Concat2 returns a [iter.Seq2] that yields all pairs from the provided [iter.Seq2]s in order.
func Concat2[I ~iterSeq2[K, V], K, V any](i ...I) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for _, s := range i {
			for k, v := range s {
				if !yield(k, v) {
					return
				}
			}
		}
	}
}

// Count2 returns the number of pairs in the provided [iter.Seq2].
func Count2[I ~iterSeq2[K, V], K, V any](i I) (o int) {
	for range i {
		o++
	}
	return
}

// Filter2 returns an [iter.Seq2] over all of the pairs for which the provided function returned true.
func Filter2[I ~iterSeq2[K, V], K, V any](i I, f func(K, V) bool) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for k, v := range i {
			if f(k, v) && !yield(k, v) {
				return
			}
		}
	}
}

// Find2 returns the first pair of the provided [iter.Seq2] that satisfies the provided function, and true.
// If no pairs satisfy the function, it returns the zero value and false.
func Find2[I ~iterSeq2[K, V], K, V any](i I, f func(K, V) bool) (k K, v V, ok bool) {
	for k, v := range i {
		if f(k, v) {
			return k, v, true
		}
	}
	return
}

// InsertMap inserts the pairs of the provided [iter.Seq2] into the provided map and returns it.
func InsertMap[O ~map[K]V, K comparable, V any](o O, i iter.Seq2[K, V]) O {
	maps.Insert(o, i)
	return o
}

// Invert returns an [iter.Seq2] over all inverted pairs.
func Invert[I ~iterSeq2[K, V], K, V any](i I) iter.Seq2[V, K] {
	return func(yield func(V, K) bool) {
		for k, v := range i {
			if !yield(v, k) {
				return
			}
		}
	}
}

// Reduce2 runs the provided function for each pair of the provided [iter.Seq2], accumulating the result in `O`.
func Reduce2[I ~iterSeq2[K, V], O, K, V any](i I, f func(O, K, V) O) (o O) {
	for k, v := range i {
		o = f(o, k, v)
	}
	return
}

// ToSeq uses the provided function to transform the pairs of the provided [iter.Seq2] into a [iter.Seq].
func ToSeq[I ~iterSeq2[K, V], K, V, E any](i I, f func(K, V) E) iter.Seq[E] {
	return func(yield func(E) bool) {
		for k, v := range i {
			if !yield(f(k, v)) {
				return
			}
		}
	}
}

// Transform2 uses the provided function to transform the pairs of the provided [iter.Seq2] into a new [iter.Seq2].
func Transform2[I ~iterSeq2[K1, V1], K1, K2, V1, V2 any](i I, f func(K1, V1) (K2, V2)) iter.Seq2[K2, V2] {
	return func(yield func(K2, V2) bool) {
		for k, v := range i {
			if !yield(f(k, v)) {
				return
			}
		}
	}
}

// TransformErr2 uses the provided function to transform the pairs of the provided [iter.Seq2] into a new [iter.Seq2].
// When the first non-nil error is encountered this function will cease iteration and set errOut to the error value.
func TransformErr2[I ~iterSeq2[K1, V1], K1, K2, V1, V2 any](i I, f func(K1, V1) (K2, V2, error), errOut *error) iter.Seq2[K2, V2] {
	return func(yield func(K2, V2) bool) {
		for ki, vi := range i {
			ko, vo, err := f(ki, vi)
			if err != nil {
				*errOut = err
				return
			}
			if !yield(ko, vo) {
				return
			}
		}
	}
}
