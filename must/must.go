package must

// Ok panics if the provided error is non-nil.
func Ok(err error) {
	if err != nil {
		panic(err)
	}
}

// Do returns the provided T, or panics if the error is non-nil.
func Do[T any](t T, err error) T {
	Ok(err)
	return t
}

// ZeroErr returns the provided value, or the zero value of the provided type if the error is non-nil.
func ZeroErr[T any](t T, err error) (zero T, _ error) {
	if err != nil {
		return zero, err
	}
	return t, nil
}
