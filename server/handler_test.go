package server

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestContainsAll(t *testing.T) {
	required := map[string][]string{
		"List1": {
			"Elem1",
			"Elem2",
		},
	}
	actual := map[string][]string{
		"List1": {
			"Elem1",
			"Elem2",
			"Elem3",
		},
		"List2": {
			"Elem1",
			"Elem2",
		},
	}
	require.True(
		t,
		containsAll(required, actual),
		"containsAll should return true if required elems are in the actual map")

	actual = map[string][]string{
		"List1": {
			"Elem1",
			"Elem2",
		},
	}
	required = map[string][]string{
		"List1": {
			"Elem1",
			"Elem2",
		},
		"List2": {
			"Elem1",
			"Elem2",
		},
	}
	require.False(
		t,
		containsAll(required, actual),
		"containsAll should return false if not all required elem are in the actual map")

	actual = map[string][]string{
		"List1": {
			"Elem1",
		},
	}
	required = map[string][]string{
		"List1": {
			"Elem1",
			"Elem2",
		},
	}
	require.False(
		t,
		containsAll(required, actual),
		"containsAll should return false if not all required elem are in the actual map")
}
