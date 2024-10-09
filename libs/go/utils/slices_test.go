package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChunks(t *testing.T) {
	t.Run("empty slice", func(t *testing.T) {
		chunks := SplitToChunks([]string{}, 5)

		assert.Len(t, chunks, 0)
	})

	t.Run("nil slice", func(t *testing.T) {
		chunks := SplitToChunks[string](nil, 5)

		assert.Len(t, chunks, 0)
	})

	t.Run("simple slice", func(t *testing.T) {
		arr := []string{"one", "two", "three", "four"}
		chunks := SplitToChunks(arr, 2)

		assert.Len(t, chunks, 2)
		assert.Len(t, chunks[0], 2)
		assert.Len(t, chunks[1], 2)

		assert.Equal(t, []string{"one", "two"}, chunks[0])
		assert.Equal(t, []string{"three", "four"}, chunks[1])
	})

	t.Run("last chunk with less elements", func(t *testing.T) {
		arr := []string{"one", "two", "three", "four", "five"}
		chunks := SplitToChunks(arr, 2)

		assert.Len(t, chunks, 3)
		assert.Equal(t, []string{"one", "two"}, chunks[0])
		assert.Equal(t, []string{"three", "four"}, chunks[1])
		assert.Equal(t, []string{"five"}, chunks[2])
	})

	t.Run("chunks of chunks", func(t *testing.T) {
		arr := make([]int, 29)
		for i := 0; i < 29; i++ {
			arr[i] = i
		}

		c1 := SplitToChunks(arr, 10)
		assert.Len(t, c1, 3)
		assert.Len(t, c1[0], 10)
		assert.Len(t, c1[2], 9)

		c2 := SplitToChunks(c1[0], 5)
		assert.Len(t, c2, 2)

		c3 := SplitToChunks(c1[2], 5)
		assert.Len(t, c3, 2)
		assert.Len(t, c3[0], 5)
		assert.Len(t, c3[1], 4)
	})

	t.Run("big chunk", func(t *testing.T) {
		arr := make([]int, 140)
		for i := 0; i < 140; i++ {
			arr[i] = i
		}

		c1 := SplitToChunks(arr, 100)
		assert.Len(t, c1, 2)
		assert.Len(t, c1[0], 100)
		assert.Len(t, c1[1], 40)
	})
}
