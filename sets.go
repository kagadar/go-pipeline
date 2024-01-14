package pipeline

import "maps"

// Difference returns the elements of `a` that are not present in `b` as a new map.
func Difference[I ~map[K]V, K comparable, V any](a, b I) I {
	o := I{}
	for k, v := range a {
		if _, ok := b[k]; !ok {
			o[k] = v
		}
	}
	return o
}

// Disjoint returns the disjoint of `a` and `b` as a new map.
func Disjoint[I ~map[K]V, K comparable, V any](a, b I) I {
	o := Difference(a, b)
	maps.Copy(o, Difference(b, a))
	return o
}

// Intersect returns the keys present in the intersection of `a` and `b`.
func Intersect[O []K, I ~map[K]V, K comparable, V any](a, b I) O {
	var s, l I
	if s, l = a, b; len(b) < len(a) {
		s, l = b, a
	}
	var o O
	for k := range s {
		if _, ok := l[k]; ok {
			o = append(o, k)
		}
	}
	return o
}

// Subset returns whether `a` is a subset of `b`.
func Subset[I ~map[K]V, K comparable, V any](a, b I) bool {
	for k := range a {
		if _, ok := b[k]; !ok {
			return false
		}
	}
	return true
}

// Union returns the union of `a` and `b` as a new map. Elements present in both will use the value from `b`.
func Union[I ~map[K]V, K comparable, V any](a, b I) I {
	o := maps.Clone(a)
	maps.Copy(o, b)
	return o
}
