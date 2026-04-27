package api

import (
	"iter"

	"github.com/kagadar/go-pipeline/slices"
)

// PaginateTokenSeq repeatedly calls the provided function, feeding in the token returned by the previous iteration, until it returns no results.
// The zero value of the token will be provided to the function on the first iteration.
//
// The function is intended to support pagination where the response contains a token to get the next results.
func PaginateTokenSeq[S ~[]E, E, T any](nextPage func(T) (S, T, error), errOut *error) iter.Seq[E] {
	return func(yield func(E) bool) {
		var token T
		for {
			s, nextToken, err := nextPage(token)
			if err != nil {
				*errOut = err
				return
			}
			var l int
			for _, e := range s {
				if !yield(e) {
					return
				}
				l++
			}
			if l == 0 {
				return
			}
			token = nextToken
		}
	}
}

// Paginate158Seq repeatedly calls the provided function, feeding in the token returned by the previous iteration, until the token is empty.
// The first iteration will be provided with an empty string for the token.
//
// The function is intended to support pagination that adheres to https://aip.dev/158.
func Paginate158Seq[S ~[]E, E any](nextPage func(string) (S, string, error), errOut *error) iter.Seq[E] {
	return func(yield func(E) bool) {
		var token string
		for {
			s, nextToken, err := nextPage(token)
			if err != nil {
				*errOut = err
				return
			}
			for _, e := range s {
				if !yield(e) {
					return
				}
			}
			if nextToken == "" {
				return
			}
			token = nextToken
		}
	}
}

// PaginateSeq repeatedly calls the provided function, taking the last value returned as an input, until it returns no results.
// The zero value of the returned type will be provided to the function on the first iteration.
//
// This function is intended to support pagination where the next page token is, or can be derived from, the last element returned.
func PaginateSeq[S ~[]E, E any](f func(E) (S, error), errOut *error) iter.Seq[E] {
	return PaginateTokenSeq(func(e E) (S, E, error) {
		s, err := f(e)
		return s, slices.Last(s), err
	}, errOut)
}
