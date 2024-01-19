package maps

import (
	"maps"
)

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

// Intersect returns the intersection of the provided maps.
// The values will be a slice of the values for that key, in order of the provided maps.
func Intersect[O map[K][]V, I ~map[K]V, K comparable, V any](i ...I) O {
	if len(i) == 0 {
		return nil
	}
	sm := i[0]
	sl := len(sm)
	for idx := 1; idx < len(i); idx++ {
		if l := len(i[idx]); l < sl {
			sm = i[idx]
			sl = l
		}
	}
	o := O{}
INPUT:
	for k := range sm {
		var vs []V
		for _, im := range i {
			if v, ok := im[k]; ok {
				vs = append(vs, v)
			} else {
				continue INPUT
			}
		}
		o[k] = vs
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

// Union returns the union of provided maps.
// The values will be a slice of the values for that key, in order of the provided maps.
func Union[O map[K][]V, I ~map[K]V, K comparable, V any](i ...I) O {
	o := O{}
	for _, m := range i {
		for k, v := range m {
			o[k] = append(o[k], v)
		}
	}
	return o
}
