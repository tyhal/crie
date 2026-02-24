package runner

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_flatten(t *testing.T) {
	m := map[string]int{"b": 1, "a": 2, "c": 3}
	flattened := flatten(m)

	slices.SortFunc(flattened, flattened.cmpV(false /*asc*/))
	expected := flattenedKVs[string, int]{
		{"b", 1},
		{"a", 2},
		{"c", 3},
	}
	assert.Equal(t, expected, flattened)

	slices.SortFunc(flattened, flattened.cmpV(true /*desc*/))
	expected = flattenedKVs[string, int]{
		{"c", 3},
		{"a", 2},
		{"b", 1},
	}
	assert.Equal(t, expected, flattened)

	slices.SortFunc(flattened, flattened.cmpK(false /*asc*/))
	expected = flattenedKVs[string, int]{
		{"a", 2},
		{"b", 1},
		{"c", 3},
	}
	assert.Equal(t, expected, flattened)
}
