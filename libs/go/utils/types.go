package utils

func Ref[T any](t T) *T {
	return &t
}

func Ok(err error) {
	if err != nil {
		panic(err)
	}
}

func Unwrap[T any](t T, err error) T {
	if err != nil {
		panic(err)
	}

	return t
}

func UnwrapError[T any](t T, err error) error {
	return err
}

func ValueOrNil[T comparable](t T) *T {
	var d T
	if t == d {
		return nil
	}

	return &t
}

func FirstNotEmpty[T comparable](values ...T) T {
	var e T

	for _, v := range values {
		if v != e {
			return v
		}
	}

	return e
}

func ValueOrDefault[T comparable](t *T) T {
	if t == nil {
		var d T
		return d
	}

	return *t
}

func CountNotNils(values ...interface{}) int {
	var count int

	for _, v := range values {
		if v != nil {
			count++
		}
	}

	return count
}
