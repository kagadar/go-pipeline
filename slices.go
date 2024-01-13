package pipeline

// Last returns the last element of the provided slice, or the zero value of that type if the slice is empty.
func Last[I ~[]E, E any](i I) (e E) {
	if l := len(i); l > 0 {
		e = i[l-1]
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

// SliceToMap uses the provided function to transform the elements of the provided slice into a map.
func SliceToMap[O map[K]V, I ~[]E, K comparable, V, E any](i I, f func(E) (K, V)) (o O) {
	o = O{}
	for _, e := range i {
		k, v := f(e)
		o[k] = v
	}
	return
}

// TransformSlice uses the provided function to transform the elements of the provided slice into a new slice.
// The returned slice does not have a stable order.
func TransformSlice[O []E2, I ~[]E1, E1, E2 any](i I, f func(E1) E2) (o O) {
	for _, e := range i {
		o = append(o, f(e))
	}
	return
}
