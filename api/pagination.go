package api

import "github.com/kagadar/go-pipeline/slices"

// PaginateToken repeatedly calls the provided function, taking a token provided by the previous iteration, until it returns no results.
// The zero value of the token will be provided to the function on the first iteration.
//
// The function is intended to support pagination compatible with aip.dev/158, where the response contains a token to get the next results.
func PaginateToken[O ~[]E, E, T any](fn func(T) (O, T, error)) (o O, err error) {
	var token T
	for {
		var s []E
		s, token, err = fn(token)
		if err != nil {
			return nil, err
		}
		if len(s) == 0 {
			return o, nil
		}
		o = append(o, s...)
	}
}

// Paginate repeatedly calls the provided function, taking the last value returned as an input, until it returns no results.
// The zero value of the returned type will be provided to the function on the first iteration.
//
// This function is intended to support pagination where the next page token is, or can be derived from, the last element returned.
func Paginate[O ~[]E, E any](fn func(E) (O, error)) (O, error) {
	return PaginateToken(func(e E) (O, E, error) {
		s, err := fn(e)
		return s, slices.Last(s), err
	})
}