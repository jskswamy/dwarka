package strings

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestContains(t *testing.T) {
	t.Run("it should return true if found", func(t *testing.T) {
		items := []string{"one", "two", "three"}

		contains := Contains(items, "one")

		assert.True(t, contains)
	})

	t.Run("it should return false if not found", func(t *testing.T) {
		items := []string{"one", "two", "three"}

		contains := Contains(items, "four")

		assert.False(t, contains)
	})
}
