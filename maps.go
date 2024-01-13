package pipeline

import (
	"cmp"
	"slices"
)

// Keys returns the keys of the provided map as a slice in no particular order.
func Keys[O []K, I ~map[K]V, K comparable, V any](i I) O {
	return MapToSlice(i, func(k K, _ V) K { return k })
}

// Values returns the values of the provided map as a slice in no particular order.
func Values[O []V, I ~map[K]V, K comparable, V any](i I) O {
	return MapToSlice(i, func(_ K, v V) V { return v })
}

// MapToSlice uses the provided function to transform the key-value pairs of the provided map into a slice.
func MapToSlice[O []E, I ~map[K]V, K comparable, V, E any](i I, f func(K, V) E) (o O) {
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

// TransformMap uses the provided function to transform the key-value pairs of the provided map into a new map.
func TransformMap[O map[K2]V2, I ~map[K1]V1, K1, K2 comparable, V1, V2 any](i I, f func(K1, V1) (K2, V2)) (o O) {
	o = make(O, len(i))
	for k, v := range i {
		k, v := f(k, v)
		o[k] = v
	}
	return
}

// RangeKeys runs the provided function once for each key-value pair in the provided map, in order of the provided keys.
func RangeKeys[IM ~map[K]V, IK ~[]K, K comparable, V any](i IM, keys IK, fn func(K, V) error) error {
	for _, k := range keys {
		v := i[k]
		if err := fn(k, v); err != nil {
			return err
		}
	}
	return nil
}

// SortedRange runs the provided function once for each key-value pair in the provided map, in ascending order of keys.
func SortedRange[I ~map[K]V, K cmp.Ordered, V any](i I, fn func(K, V) error) error {
	s := Keys(i)
	slices.Sort(s)
	return RangeKeys(i, s, fn)
}

// SortedRangeFunc runs the provided range function once for each key-value pair in the provided map, in ascending order of keys as determined by the provided sort function.
func SortedRangeFunc[I ~map[K]V, K comparable, V any](i I, sortFn func(x, y K) int, rangeFn func(K, V) error) error {
	s := Keys(i)
	slices.SortFunc(s, sortFn)
	return RangeKeys(i, s, rangeFn)
}
