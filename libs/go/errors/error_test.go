package errors

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateError(t *testing.T) {
	type test struct {
		name     string
		code     *ErrorCode
		message  string
		cause    error
		metadata Metadata
	}

	tests := []test{
		{
			name:     "error with message only",
			code:     Define("basic_error"),
			message:  "custom message",
			metadata: make(Metadata),
		},
		{
			name:    "error with message and metadata",
			code:    Define("error_with_metadata"),
			message: "custom message",
			metadata: Metadata{
				"str": "value",
				"num": 123,
			},
		},
		{
			name:     "wrap error with cause",
			code:     Define("error_with_metadata"),
			message:  "custom message",
			cause:    errors.New("raw error"),
			metadata: make(Metadata),
		},
		{
			name:    "error with message, metadata and cause",
			code:    Define("error_with_metadata"),
			message: "custom message",
			cause:   errors.New("raw error"),
			metadata: Metadata{
				"str": "value",
				"num": 123,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			m := NewMetadata()
			for k, v := range test.metadata {
				m = m.And(k, v)
			}

			var err *Error
			if test.cause != nil {
				err = Wrap(test.code, test.cause, test.message, m)
			} else {
				err = New(test.code, test.message, m)
			}

			assert.Equal(t, test.code, err.code)
			assert.Equal(t, test.message, err.message)
			assert.Equal(t, test.cause, err.cause)
			assert.Equal(t, test.metadata, err.metadata)
		})
	}
}

func TestIsError(t *testing.T) {
	code1 := Define("code1")
	code2 := Define("code2")

	err1 := New(code1, "message")
	err2 := New(code2, "message")
	err3 := Wrap(code1, err2, "wrapped")

	assert.True(t, err1.Is(code1))

	assert.True(t, err1.Is(err1))

	assert.False(t, err1.Is(code2))

	assert.False(t, err1.Is(err2))

	assert.True(t, err3.Is(code1))
	assert.True(t, err3.Is(code2), "recursive")
	assert.True(t, err3.Is(err2), "recursive")
	assert.True(t, err3.Is(err1), "recursive and same code")
}

func TestMarshalError(t *testing.T) {
	rawErr := errors.New("raw error")
	customErr := Wrap(
		Define("code"),
		rawErr,
		"custom message",
		NewMetadata().And("key", "value"),
	)

	errJson, err := json.Marshal(customErr)
	assert.NoError(t, err)
	assert.Equal(
		t,
		[]byte(`{"cause":{"message":"raw error"},"code":"code","message":"custom message","metadata":{"key":"value"}}`),
		errJson,
	)
}
