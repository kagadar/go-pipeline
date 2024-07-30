package api

import "github.com/kagadar/go-pipeline/slices"

// PaginateToken repeatedly calls the provided function, feeding in the token returned by the previous iteration, until it returns no results.
// The zero value of the token will be provided to the function on the first iteration.
//
// The function is intended to support pagination where the response contains a token to get the next results.
func PaginateToken[O ~[]E, E, T any](nextPage func(T) (O, T, error)) (o O, err error) {
	var token T
	for {
		var s []E
		s, token, err = nextPage(token)
		if err != nil {
			return nil, err
		}
		if len(s) == 0 {
			return o, nil
		}
		o = append(o, s...)
	}
}

// Paginate158 repeatedly calls the provided function, feeding in the token returned by the previous iteration, until the token is empty.
// The first iteration will be provided with an empty string for the token.
//
// The function is intended to support pagination that adheres to https://aip.dev/158.
func Paginate158[O ~[]E, E any](nextPage func(string) (O, string, error)) (o O, err error) {
	var token string
	for {
		var es []E
		es, token, err = nextPage(token)
		if err != nil {
			return nil, err
		}
		o = append(o, es...)
		if token == "" {
			return o, nil
		}
	}
}

// Paginate repeatedly calls the provided function, taking the last value returned as an input, until it returns no results.
// The zero value of the returned type will be provided to the function on the first iteration.
//
// This function is intended to support pagination where the next page token is, or can be derived from, the last element returned.
func Paginate[O ~[]E, E any](f func(E) (O, error)) (O, error) {
	return PaginateToken(func(e E) (O, E, error) {
		s, err := f(e)
		return s, slices.Last(s), err
	})
}
