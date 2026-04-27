package predicates

import "cmp"

type Integer interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

type Parity int

const (
	Even Parity = iota
	Odd
)

func IsEven[T Integer](t T) bool { return t%2 == 0 }

func FindParity[T Integer](t T) Parity {
	if IsEven(t) {
		return Even
	}
	return Odd
}

func IsZero[T comparable](t T) bool {
	var zero T
	return t == zero
}

func Keys[K, V, O any](f func(K) O) func(K, V) O {
	return func(k K, _ V) O { return f(k) }
}

func Values[K, V, O any](f func(V) O) func(K, V) O {
	return func(_ K, v V) O { return f(v) }
}

func ByMapValue[I ~map[K]V, K comparable, V cmp.Ordered](i I) func(a, b K) int {
	return func(a, b K) int {
		return cmp.Compare(i[a], i[b])
	}
}

func ByMapValueFunc[I ~map[K]V, K comparable, V any](i I, f func(a, b V) int) func(a, b K) int {
	return func(a, b K) int {
		return f(i[a], i[b])
	}
}
