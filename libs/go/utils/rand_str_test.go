package utils

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRandFromString(t *testing.T) {
	rand1 := NewRandFromString("test")
	rand2 := NewRandFromString("test")
	assert.Equal(t, rand1.Int(), rand2.Int())
	assert.Equal(t, rand1.Float64(), rand2.Float64())

	rand3 := NewRandFromString("test")
	rand4 := NewRandFromString("test2")
	assert.NotEqual(t, rand3.Int(), rand4.Int())
	assert.NotEqual(t, rand3.Float64(), rand4.Float64())
}

func TestPercentageFromString(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		n := PercentageFromString("hello")
		assert.Equal(t, 75.342, n)

		n = PercentageFromString("test")
		assert.Equal(t, 42.533, n)

		n = PercentageFromString("")
		assert.Equal(t, 61.652, n)
	})

	t.Run("between 0 and 100", func(t *testing.T) {
		for i := 0; i < 250; i++ {
			str := fmt.Sprintf("test-%d", i)
			n := PercentageFromString(str)
			assert.Less(t, n, 100.0, "str: %s - n: %f", str, n)
			assert.Greater(t, n, 0.0, "str: %s - n: %f", str, n)
		}
	})
}
