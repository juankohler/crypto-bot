package errors

import (
	"encoding/json"
	"fmt"
	"strings"
)

/** ErrorCode */
type ErrorCode struct {
	code string
}

func Define(code string) *ErrorCode {
	if len(code) == 0 {
		panic("empty error code")
	}

	return &ErrorCode{code}
}

func (c *ErrorCode) String() string {
	return c.code
}

func (c *ErrorCode) Error() string {
	return c.code
}

/** Error */
type Error struct {
	code     *ErrorCode
	message  string
	cause    error
	metadata Metadata
}

func New(code *ErrorCode, message string, metadata ...Metadata) *Error {
	m := NewMetadata()
	for _, metadata := range metadata {
		m = m.Merge(metadata)
	}

	return &Error{
		code:     code,
		message:  message,
		metadata: m,
	}
}

func Wrap(code *ErrorCode, cause error, message string, metadata ...Metadata) *Error {
	err := New(code, message, metadata...)
	err.cause = cause

	return err
}

func (err *Error) Code() *ErrorCode {
	return err.code
}

func (err *Error) Message() string {
	return err.message
}

func (err *Error) Cause() error {
	return err.cause
}

func (err *Error) Metadata() Metadata {
	return err.metadata
}

func (err *Error) Error() string {
	var str string

	if err.code == nil || err.cause == nil {
		return err.message
	}

	str = fmt.Sprintf("%s: %s", err.code.code, err.message)

	if err.cause != nil {
		str += fmt.Sprintf(" (%s)", err.cause.Error())
	}

	if len(err.metadata) > 0 {
		metadataStr := make([]string, 0, len(err.metadata))
		for k, v := range err.metadata {
			metadataStr = append(metadataStr, fmt.Sprintf("[%s = %v]", k, v))
		}

		str += fmt.Sprintf(" %s", strings.Join(metadataStr, ", "))
	}

	return str
}

func (err *Error) Unwrap() error {
	return err.cause
}

func (err *Error) Is(target error) bool {
	if target == nil {
		return false
	}

	if t, ok := target.(*Error); ok {
		target = t.code
	}

	if t, ok := target.(*ErrorCode); ok {
		if err.code == t {
			return true
		}

		if err.cause != nil {
			if cause, ok := err.cause.(*Error); ok {
				return cause.Is(t)
			}
		}
	}

	return false
}

// Serialization
func (err *Error) MarshalJSON() ([]byte, error) {
	var cause interface{}
	if err.cause != nil {
		switch err := err.cause.(type) {
		case *Error:
			cause = err
		case error:
			cause = map[string]interface{}{
				"message": err.Error(),
			}
		}
	}

	return json.Marshal(map[string]interface{}{
		"code":     err.code.code,
		"message":  err.message,
		"cause":    cause,
		"metadata": err.metadata,
	})
}
