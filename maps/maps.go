package maps

import (
	"cmp"
	"maps"
	"slices"

	"github.com/kagadar/go-pipeline/must"
	"github.com/kagadar/go-pipeline/predicates"
	"github.com/kagadar/go-pipeline/seq"
)

// Keys returns the keys of the provided map as a slice in no particular order.
func Keys[I ~map[K]V, K comparable, V any](i I) []K {
	if len(i) == 0 {
		return nil
	}
	return seq.CollectSlice[[]K](len(i), maps.Keys(i))
}

// Values returns the values of the provided map as a slice in no particular order.
func Values[I ~map[K]V, K comparable, V any](i I) []V {
	if len(i) == 0 {
		return nil
	}
	return seq.CollectSlice[[]V](len(i), maps.Values(i))
}

// All returns whether all key-value pairs of the provided map satisfy the provided function.
func All[I ~map[K]V, K comparable, V any](i I, f func(K, V) bool) bool {
	return seq.All2(maps.All(i), f)
}

// Any returns whether any key-value pair of the provided map satisfies the provided function.
func Any[I ~map[K]V, K comparable, V any](i I, f func(K, V) bool) bool {
	return seq.Any2(maps.All(i), f)
}

// Filter returns a new map containing all of the key-value pairs for which the provided function returned true.
// The returned map will be preallocated to len(i). To avoid pre-allocation, use [seq.Filter2].
func Filter[I ~map[K]V, K comparable, V any](i I, f func(K, V) bool) I {
	return seq.CollectMap[I](len(i), seq.Filter2(maps.All(i), f))
}

// Invert returns a new map with all key-value pairs inverted.
func Invert[I ~map[K]V, K, V comparable](i I) map[V]K {
	return seq.CollectMap[map[V]K](len(i), seq.Invert(maps.All(i)))

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

// Partition returns two new maps of the key-value pairs of the provided map, splitting them based on whether the provided func returned true.
// The `pass` and `fail` maps will be pre-allocated to len(i)/2.
func Partition[I ~map[K]V, K comparable, V any](i I, f func(K, V) bool) (pass, fail I) {
	pass, fail = make(I, len(i)/2), make(I, len(i)/2)
	for k, v := range i {
		if f(k, v) {
			pass[k] = v
		} else {
			fail[k] = v
		}
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
func Reduce[I ~map[K]V, K comparable, V, O any](i I, f func(O, K, V) O) O {
	return seq.Reduce2(maps.All(i), f)
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

// ToSlice uses the provided function to transform the key-value pairs of the provided map into a slice.
func ToSlice[I ~map[K]V, K comparable, V, E any](i I, f func(K, V) E) []E {
	if len(i) == 0 {
		return nil
	}
	return seq.CollectSlice[[]E](len(i), seq.ToSeq(maps.All(i), f))
}

// Transform uses the provided function to transform the key-value pairs of the provided map into a new map.
func Transform[I ~map[K1]V1, K1, K2 comparable, V1, V2 any](i I, f func(K1, V1) (K2, V2)) map[K2]V2 {
	return seq.CollectMap[map[K2]V2](len(i), seq.Transform2(maps.All(i), f))
}

// TransformErr uses the provided function to transform the key-value pairs of the provided map into a new map.
// When the first non-nil error is encountered this function will return a nil map and the error without transforming the remaining elements.
func TransformErr[I ~map[K1]V1, K1, K2 comparable, V1, V2 any](i I, f func(K1, V1) (K2, V2, error)) (o map[K2]V2, err error) {
	return must.ZeroErr(seq.CollectMap[map[K2]V2](len(i), seq.TransformErr2(maps.All(i), f, &err)), err)
}

// ValueSortedRange runs the provided range function once for each key-value pair in the provided map, in ascending order of values.
func ValueSortedRange[I ~map[K]V, K comparable, V cmp.Ordered](i I, f func(K, V) error) error {
	keys := Keys(i)
	slices.SortFunc(keys, predicates.ByMapValue(i))
	return Range(i, keys, f)
}

// ValueSortedRangeFunc runs the provided range function once for each key-value pair in the provided map, in ascending order of values as determined by the provided sort function.
func ValueSortedRangeFunc[I ~map[K]V, K comparable, V any](i I, sortF func(x, y V) int, rangeF func(K, V) error) error {
	keys := Keys(i)
	slices.SortFunc(keys, predicates.ByMapValueFunc(i, sortF))
	return Range(i, keys, rangeF)
}
