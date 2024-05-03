package chains

// MapLink wraps a map and error together to allow functions which take and
// return links to be chained together.
type MapLink[K comparable, V any] func() (map[K]V, error)

// NewMapLink is a convenience function to create a new MapLink which wraps a
// map and has no error set.
func NewMapLink[K comparable, V any](m map[K]V) MapLink[K, V] {
	return func() (map[K]V, error) {return m, nil}
}
