package pipeline

func Flatten[I ~[]O, O ~[]E, E any](i I) (o O) {
	for _, e := range i {
		o = append(o, e...)
	}
	return
}

func SliceToMap[I ~[]E, O ~map[K]V, K comparable, V, E any](i I, f func(E) (K, V)) (o O) {
	for _, e := range i {
		k, v := f(e)
		o[k] = v
	}
	return
}

func TransformSlice[I ~[]E1, O []E2, E1, E2 any](i I, f func(E1) E2) (o O) {
	for _, e := range i {
		o = append(o, f(e))
	}
	return
}
