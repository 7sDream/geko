package geko_test

import (
	"reflect"
	"testing"

	"github.com/7sDream/geko"
)

func TestPairList_New(t *testing.T) {
	kol := geko.NewPairList[string, int]()

	if kol.List != nil {
		t.Fatalf("NewPairList inner slice is not nil")
	}

	list := []geko.Pair[string, int]{
		{"one", 1},
		{"two", 2},
		{"three", 3},
	}
	kol2 := geko.NewPairListFrom(list)

	if !reflect.DeepEqual(kol2.List, list) {
		t.Fatalf("NewPairList doesn't store origin slice")
	}
}

func TestPairList_NewWithCapacity(t *testing.T) {
	kol := geko.NewPairListWithCapacity[string, int](12)

	if cap(kol.List) != 12 {
		t.Fatalf("NewPairListWithCapacity inner slice does not have correct capacity")
	}
}
