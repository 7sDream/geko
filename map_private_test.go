package geko

import "testing"

func TestMap_NewWithCapacity(t *testing.T) {
	m := NewMapWithCapacity[string, int](20)

	if cap(m.order) != 20 {
		t.Fatalf("NewMapWithCapacity does not init with capacity")
	}
}
