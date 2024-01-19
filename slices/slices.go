package slices

import "slices"

// Filter returns a new slice containing all of the elements for which the provided function returned true, in order.
func Filter[I ~[]E, E any](i I, f func(E) bool) (o I) {
	for _, e := range i {
		if f(e) {
			o = append(o, e)
		}
	}
	return
}

// Flatten returns a new slice containing all of the elements, in order of the provided slice of slices.
func Flatten[O []E, I ~[]O, E any](i I) (o O) {
	for _, e := range i {
		o = append(o, e...)
	}
	return
}

// Last returns the last element of the provided slice, or the zero value of that type if the slice is empty.
func Last[I ~[]E, E any](i I) (e E) {
	if l := len(i); l > 0 {
		e = i[l-1]
	}
	return
}

// Reduce runs the provided function once for each element, in order, accumulating the result in `O`.
func Reduce[O any, I ~[]E, E any](i I, f func(O, E) O) (o O) {
	for _, e := range i {
		o = f(o, e)
	}
	return
}

// ToMap uses the provided function to transform the elements of the provided slice into a map.
func ToMap[O map[K]V, I ~[]E, K comparable, V, E any](i I, f func(E) (K, V)) (o O) {
	o = O{}
	for _, e := range i {
		k, v := f(e)
		o[k] = v
	}
	return
}

// Transform uses the provided function to transform the elements of the provided slice into a new slice.
// The returned slice does not have a stable order.
func Transform[O []E2, I ~[]E1, E1, E2 any](i I, f func(E1) E2) (o O) {
	for _, e := range i {
		o = append(o, f(e))
	}
	return
}

// Zip combines the element at each index from each slice into a new slice of the input type.
// These slices are returned in order, from 0 to the length of the shortest input slice.
func Zip[O []I, I ~[]E, E any](i ...I) O {
	o := make(O, slices.Min(Transform(i, func(i I) int { return len(i) })))
	for idx := 0; idx < len(o); idx++ {
		o[idx] = make(I, len(i))
		for ii, is := range i {
			o[idx][ii] = is[idx]
		}
	}
	return o
}
