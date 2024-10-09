package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFirstNotEmpty(t *testing.T) {
	assert.Equal(t, "a", FirstNotEmpty("a", "b", "c"))
	assert.Equal(t, "b", FirstNotEmpty("", "b", "c"))
	assert.Equal(t, "c", FirstNotEmpty("", "", "c"))
	assert.Equal(t, "", FirstNotEmpty("", "", ""))
}
