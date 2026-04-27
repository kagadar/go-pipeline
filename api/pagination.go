package api

import (
	"slices"

	"github.com/kagadar/go-pipeline/must"
)

// PaginateToken repeatedly calls the provided function, feeding in the token returned by the previous iteration, until it returns no results.
// The zero value of the token will be provided to the function on the first iteration.
//
// The function is intended to support pagination where the response contains a token to get the next results.
func PaginateToken[O ~[]E, E, T any](nextPage func(T) (O, T, error)) (_ O, err error) {
	return must.ZeroErr(slices.Collect(PaginateTokenSeq(nextPage, &err)), err)
}

// Paginate158 repeatedly calls the provided function, feeding in the token returned by the previous iteration, until the token is empty.
// The first iteration will be provided with an empty string for the token.
//
// The function is intended to support pagination that adheres to https://aip.dev/158.
func Paginate158[O ~[]E, E any](nextPage func(string) (O, string, error)) (_ O, err error) {
	return must.ZeroErr(slices.Collect(Paginate158Seq(nextPage, &err)), err)
}

// Paginate repeatedly calls the provided function, taking the last value returned as an input, until it returns no results.
// The zero value of the returned type will be provided to the function on the first iteration.
//
// This function is intended to support pagination where the next page token is, or can be derived from, the last element returned.
func Paginate[O ~[]E, E any](f func(E) (O, error)) (_ O, err error) {
	return must.ZeroErr(slices.Collect(PaginateSeq(f, &err)), err)
}
