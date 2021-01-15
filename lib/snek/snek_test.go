package snek

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDesnekRecursive(t *testing.T) {
	t.Run("Nested Snake Case", func(t *testing.T) {
		// Given
		input := map[string]interface{}{
			"foo_bar": map[string]interface{}{
				"bar_foo": false,
			},
			"integer_five": 5,
		}

		// When
		actualResult := DesnekKeysRecursive(input)

		// Then
		expectedResult := map[string]interface{}{
			"fooBar": map[string]interface{}{
				"barFoo": false,
			},
			"integerFive": 5,
		}
		assert.Equal(t, expectedResult, actualResult)
	})
}

func TestSnekRecursive(t *testing.T) {
	t.Run("Nested Camel Case", func(t *testing.T) {
		// Given
		input := map[string]interface{}{
			"fooBar": map[string]interface{}{
				"barFoo": false,
			},
			"integerFive": 5,
		}

		// When
		actualResult := SnekKeysRecursive(input)

		// Then
		expectedResult := map[string]interface{}{
			"foo_bar": map[string]interface{}{
				"bar_foo": false,
			},
			"integer_five": 5,
		}
		assert.Equal(t, expectedResult, actualResult)
	})
}
