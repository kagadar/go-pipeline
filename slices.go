package pipeline

func Flatten[O []E, I ~[]O, E any](i I) (o O) {
	for _, e := range i {
		o = append(o, e...)
	}
	return
}

func SliceToMap[O map[K]V, I ~[]E, K comparable, V, E any](i I, f func(E) (K, V)) (o O) {
	o = O{}
	for _, e := range i {
		k, v := f(e)
		o[k] = v
	}
	return
}

func TransformSlice[O []E2, I ~[]E1, E1, E2 any](i I, f func(E1) E2) (o O) {
	for _, e := range i {
		o = append(o, f(e))
	}
	return
}
