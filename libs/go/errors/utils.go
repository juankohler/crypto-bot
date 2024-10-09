package errors

import (
	"errors"
)

func Is(err error, target error) bool {
	if err == nil || target == nil {
		return false
	}

	if err, ok := err.(*Error); ok {
		return err.Is(target)
	}

	return errors.Is(err, target)
}
