package mapchains

import (
	"errors"
	"fmt"

	"github.com/kagadar/go-pipeline/chains"
)

var errKeyCollision = errors.New("key already in map")

func IsKeyCollisionErr(err error) bool {
	return errors.Is(err, errKeyCollision)
}

func propagateErr[K comparable, V any](err error) chains.MapLink[K, V] {
	return func() (map[K]V, error) { return nil, err }
}

// Transform uses the provided function to transform the key-value pairs of the
// provided MapLink into a new MapLink, with new key-value types.
func Transform[K1, K2 comparable, V1, V2 any](i chains.MapLink[K1, V1], f func(K1, V1) (K2, V2, error)) chains.MapLink[K2, V2] {
	in, err := i()
	if err != nil {
		return propagateErr[K2, V2](err)
	}
	o := make(map[K2]V2, len(in))
	for k, v := range in {
		k, v, err := f(k, v)
		if err != nil {
			return propagateErr[K2, V2](err)
		}
		if _, ok := o[k]; ok {
			return propagateErr[K2, V2](fmt.Errorf("%v: %w", k, errKeyCollision))
		}
		o[k] = v
	}
	return chains.NewMapLink(o)
}

// Filter returns a new map containing all of the key-value pairs for which the provided function returned true.
func Filter[K comparable, V any](i chains.MapLink[K, V], f func(K, V) (bool, error)) chains.MapLink[K, V] {
	in, err := i()
	if err != nil {
		return propagateErr[K, V](err)
	}
	o := make(map[K]V, len(in))
	for k, v := range in {
		keep, err := f(k, v)
		if err != nil {
			return propagateErr[K, V](err)
		}
		if keep {
			o[k] = v
		}
	}
	return chains.NewMapLink(o)
}
