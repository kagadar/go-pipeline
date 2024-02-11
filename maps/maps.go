package maps

import (
	"cmp"
	"slices"

	pslices "github.com/kagadar/go-pipeline/slices"
)

// Keys returns the keys of the provided map as a slice in no particular order.
func Keys[O []K, I ~map[K]V, K comparable, V any](i I) O {
	return ToSlice(i, func(k K, _ V) K { return k })
}

// Values returns the values of the provided map as a slice in no particular order.
func Values[O []V, I ~map[K]V, K comparable, V any](i I) O {
	return ToSlice(i, func(_ K, v V) V { return v })
}

// Filter returns a new map containing all of the key-value pairs for which the provided function returned true.
func Filter[I ~map[K]V, K comparable, V any](i I, f func(K, V) bool) I {
	o := I{}
	for k, v := range i {
		if f(k, v) {
			o[k] = v
		}
	}
	return o
}

// ToSlice uses the provided function to transform the key-value pairs of the provided map into a slice.
func ToSlice[O []E, I ~map[K]V, K comparable, V, E any](i I, f func(K, V) E) (o O) {
	o = make(O, 0, len(i))
	for k, v := range i {
		o = append(o, f(k, v))
	}
	return
}

// MapMapInsert inserts the provided value into the inner map of the provided map of maps.
// If the inner map is nil, it will be initialised first.
func MapMapInsert[I ~map[K]I2, I2 ~map[K2]V, K, K2 comparable, V any](i I, k K, k2 K2, v V) {
	if i2, ok := i[k]; !ok {
		i[k] = I2{k2: v}
	} else {
		i2[k2] = v
	}
}

// Transform uses the provided function to transform the key-value pairs of the provided map into a new map.
func Transform[O map[K2]V2, I ~map[K1]V1, K1, K2 comparable, V1, V2 any](i I, f func(K1, V1) (K2, V2)) (o O) {
	o = make(O, len(i))
	for k, v := range i {
		k, v := f(k, v)
		o[k] = v
	}
	return
}

// Range runs the provided function once for each key-value pair in the provided map, in order of the provided keys.
func Range[IM ~map[K]V, IK ~[]K, K comparable, V any](i IM, keys IK, f func(K, V) error) error {
	for _, k := range keys {
		v := i[k]
		if err := f(k, v); err != nil {
			return err
		}
	}
	return nil
}

// Reduce runs the provided function once for each key-value pair, accumulating the result in `O`.
func Reduce[O any, I ~map[K]V, K comparable, V any](i I, f func(O, K, V) O) (o O) {
	for k, v := range i {
		o = f(o, k, v)
	}
	return
}

// SortedRange runs the provided function once for each key-value pair in the provided map, in ascending order of keys.
func SortedRange[I ~map[K]V, K cmp.Ordered, V any](i I, f func(K, V) error) error {
	s := Keys(i)
	slices.Sort(s)
	return Range(i, s, f)
}

// SortedRangeFunc runs the provided range function once for each key-value pair in the provided map, in ascending order of keys as determined by the provided sort function.
func SortedRangeFunc[I ~map[K]V, K comparable, V any](i I, sortF func(x, y K) int, rangeF func(K, V) error) error {
	s := Keys(i)
	slices.SortFunc(s, sortF)
	return Range(i, s, rangeF)
}

// ValueSortedRange runs the provided range function once for each key-value pair in the provided map, in ascending order of values.
func ValueSortedRange[I ~map[K]V, K comparable, V cmp.Ordered](i I, f func(K, V) error) error {
	type pair struct {
		k K
		v V
	}
	pairs := ToSlice(i, func(k K, v V) pair {
		return pair{k, v}
	})
	slices.SortFunc(pairs, func(x, y pair) int {
		if x.v < y.v {
			return -1
		}
		if x.v > y.v {
			return 1
		}
		return 0
	})
	return Range(i, pslices.Transform(pairs, func(p pair) K { return p.k }), f)
}

// ValueSortedRange runs the provided range function once for each key-value pair in the provided map, in ascending order of values as determined by the provided sort function.
func ValueSortedRangeFunc[I ~map[K]V, K comparable, V any](i I, sortF func(x, y V) int, rangeF func(K, V) error) error {
	type pair struct {
		k K
		v V
	}
	pairs := ToSlice(i, func(k K, v V) pair {
		return pair{k, v}
	})
	slices.SortFunc(pairs, func(x, y pair) int {
		return sortF(x.v, y.v)
	})
	return Range(i, pslices.Transform(pairs, func(p pair) K { return p.k }), rangeF)
}
