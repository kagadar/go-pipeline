package slices

import (
	"slices"

	"github.com/kagadar/go-pipeline/must"
	"github.com/kagadar/go-pipeline/seq"
)

// All returns whether all elements of the provided slice satisfy the provided function.
func All[I ~[]E, E any](i I, f func(E) bool) bool {
	return seq.All(slices.Values(i), f)
}

// Dedupe returns a new slice containing the distinct elements of the provided slice, in order of first occurrence.
// The returned slice will be preallocated to len(i). To avoid pre-allocation, use [seq.Dedupe].
func Dedupe[I ~[]E, E comparable](i I) I {
	if len(i) == 0 {
		return nil
	}
	return seq.CollectSlice[I](len(i), seq.Dedupe(slices.Values(i)))
}

// Filter returns a new slice containing all of the elements for which the provided function returned true.
// The returned slice will be preallocated to len(i). To avoid pre-allocation, use [seq.Filter].
func Filter[I ~[]E, E any](i I, f func(E) bool) I {
	if len(i) == 0 {
		return nil
	}
	return seq.CollectSlice[I](len(i), seq.Filter(slices.Values(i), f))
}

// Flatten returns a new slice containing all of the elements in order from the provided slice of slices.
func Flatten[I ~[]S, S ~[]E, E any](i I) S {
	if len(i) == 0 {
		return nil
	}
	return slices.Concat(i...)
}

// GroupBy returns a map of the elements of the provided slice grouped by the value returned by running the key func on them.
func GroupBy[I ~[]E, K comparable, E any](i I, f func(E) K) map[K][]E {
	o := map[K][]E{}
	for _, e := range i {
		k := f(e)
		o[k] = append(o[k], e)
	}
	return o
}

// Last returns the last element of the provided slice, or the zero value of that type if the slice is empty.
func Last[I ~[]E, E any](i I) (e E) {
	if l := len(i); l > 0 {
		e = i[l-1]
	}
	return
}

// Partition returns two new slices of the elements of the provided slice, splitting them based on whether the provided func returned true.
// To avoid reallocation, the `fail` slice will be in reverse order.
func Partition[I ~[]E, E any](i I, f func(E) bool) (I, I) {
	if len(i) == 0 {
		return nil, nil
	}
	out := make(I, len(i))
	front, back := 0, len(i)-1
	for _, e := range i {
		if f(e) {
			out[front] = e
			front++
		} else {
			out[back] = e
			back--
		}
	}
	return out[:front:front], out[front:]
}

// Reduce runs the provided function for each element of the provided slice, accumulating the result in `O`.
func Reduce[I ~[]E, E, O any](i I, f func(O, E) O) O {
	return seq.Reduce(slices.Values(i), f)
}

// SkipWhile returns a new slice containing the elements of the provided slice, skipping elements from the start until provided function returns true for the first time.
// The returned slice will be preallocated to len(i). To avoid pre-allocation, use [seq.SkipWhile].
func SkipWhile[I ~[]E, E any](i I, f func(E) bool) I {
	if len(i) == 0 {
		return nil
	}
	return seq.CollectSlice[I](len(i), seq.SkipWhile(slices.Values(i), f))
}

// TakeWhile returns a new slice containing the elements of the provided slice until the provided function returns false.
// The returned slice will be preallocated to len(i). To avoid pre-allocation, use [seq.TakeWhile].
func TakeWhile[I ~[]E, E any](i I, f func(E) bool) I {
	if len(i) == 0 {
		return nil
	}
	return seq.CollectSlice[I](len(i), seq.TakeWhile(slices.Values(i), f))
}

// ToMap uses the provided function to transform the elements of the provided slice into a map.
// The returned map will be preallocated to len(i). To avoid pre-allocation, use [seq.ToSeq2].
func ToMap[I ~[]E, K comparable, V, E any](i I, f func(E) (K, V)) map[K]V {
	return seq.CollectMap[map[K]V](len(i), seq.ToSeq2(slices.Values(i), f))
}

// Transform uses the provided function to transform the elements of the provided slice into a new slice.
func Transform[I ~[]E1, E1, E2 any](i I, f func(E1) E2) []E2 {
	if len(i) == 0 {
		return nil
	}
	return seq.CollectSlice[[]E2](len(i), seq.Transform(slices.Values(i), f))
}

// TransformErr uses the provided function to transform the elements of the provided slice into a new slice.
// When the first non-nil error is encountered this function will return a nil slice and the error without transforming the remaining elements.
func TransformErr[I ~[]E1, E1, E2 any](i I, f func(E1) (E2, error)) (_ []E2, err error) {
	if len(i) == 0 {
		return nil, nil
	}
	return must.ZeroErr(seq.CollectSlice[[]E2](len(i), seq.TransformErr(slices.Values(i), f, &err)), err)
}

// Zip combines the slices into a slice of slices, where each inner slice contains the elements at the corresponding index of each input slice.
// These slices are returned in order, from 0 to the length of the shortest input slice.
func Zip[I ~[]E, E any](i ...I) []I {
	if len(i) == 0 {
		return nil
	}
	l := len(i[0])
	for _, is := range i[1:] {
		if len(is) < l {
			l = len(is)
		}
	}
	o := make([]I, l)
	for idx := range o {
		o[idx] = make(I, len(i))
		for ii, is := range i {
			o[idx][ii] = is[idx]
		}
	}
	return o
}
