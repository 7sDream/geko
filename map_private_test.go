package geko

import "testing"

func TestMap_NewWithCapacity(t *testing.T) {
	kom := NewMapWithCapacity[string, int](20)

	if cap(kom.order) != 20 {
		t.Fatalf("NewMapWithCapacity does not init with capacity")
	}
}
