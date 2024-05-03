package mapchains

import (
	"github.com/kagadar/go-pipeline/chains"
)

func propagateErr[K comparable, V any](err error) chains.MapLink[K, V] {
	return func() (map[K]V, error) { return nil, err }
}

// Transform uses the provided function to transform the key-value pairs of the provided MapLink into a new MapLink.
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
		o[k] = v
	}
	return chains.NewMapLink(o)
}
