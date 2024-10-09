package errors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUtils(t *testing.T) {
	t.Run("is", func(t *testing.T) {
		code1 := Define("code")
		code2 := Define("another-code")

		err := New(code1, "message")
		wrapped := Wrap(code2, err, "wrapped")

		// Error
		assert.True(t, Is(err, code1))
		assert.True(t, Is(wrapped, code2))
		assert.True(t, Is(wrapped, code1), "recursive")
		assert.False(t, Is(err, Define("code")), "different instances")

		// Code
		assert.True(t, Is(code1, code1), "same instance")
		assert.False(t, Is(code1, code2))
		assert.False(t, Is(code1, Define("another-code")))
		assert.False(t, Is(code1, Define("code")), "different instances")
	})
}
