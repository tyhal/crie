package crie

import (
	"reflect"
	"strings"
	"testing"
)

func TestFilter(t *testing.T) {
	list := []string{"apple", "banana", "cherry", "apricot"}

	// Test expect=true
	result := Filter(list, true, func(s string) bool { return strings.HasPrefix(s, "a") })
	expected := []string{"apple", "apricot"}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("got %v, want %v", result, expected)
	}

	// Test expect=false
	result = Filter(list, false, func(s string) bool { return strings.HasPrefix(s, "a") })
	expected = []string{"banana", "cherry"}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("got %v, want %v", result, expected)
	}

	// Test empty result
	result = Filter([]string{}, true, func(s string) bool { return true })
	if len(result) != 0 {
		t.Errorf("expected empty slice, got %v", result)
	}
}

func BenchmarkFilter(b *testing.B) {
	list := []string{"apple", "banana", "cherry", "date"}
	filterFn := func(s string) bool { return len(s) > 4 }

	for i := 0; i < b.N; i++ {
		Filter(list, true, filterFn)
	}
}
