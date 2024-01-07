package pipeline

func MapToSlice[I ~map[K]V, O ~[]E, K comparable, V, E any](i I, f func(K, V) E) (o O) {
	for k, v := range i {
		o = append(o, f(k, v))
	}
	return
}

func TransformMap[I ~map[K1]V1, O ~map[K2]V2, K1, K2 comparable, V1, V2 any](i I, f func(K1, V1) (K2, V2)) (o O) {
	for k, v := range i {
		k, v := f(k, v)
		o[k] = v
	}
	return
}
