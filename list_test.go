package geko_test

import (
	"reflect"
	"testing"

	"github.com/7sDream/geko"
)

func TestList_New(t *testing.T) {
	kol := geko.NewList[int]()

	if kol.List != nil {
		t.Fatalf("NewList inner slice is not nil")
	}

	kol2 := geko.NewListFrom([]int{1, 2, 3})

	if reflect.DeepEqual(kol2.List, []int{1, 2, 3}) {
		t.Fatalf("NewListFrom doesn't store origin slice")
	}
}

func TestList_NewWithCapacity(t *testing.T) {
	kol := geko.NewListWithCapacity[int](12)

	if cap(kol.List) != 12 {
		t.Fatalf("NewListWithCapacity inner slice does not have correct capacity")
	}
}

func TestList_Get(t *testing.T) {
	kol := geko.NewListFrom([]int{1, 2, 3})

	if !willPanic(func() {
		kol.Get(-1)
	}) {
		t.Fatalf("Get doesn't panic with negative index")
	}

	if !willPanic(func() {
		kol.Get(4)
	}) {
		t.Fatalf("Get doesn't panic with out-of-bound index")
	}

	if kol.Get(1) != 2 {
		t.Fatalf("Get doesn't return correct value")
	}
}

func TestList_Set(t *testing.T) {
	kol := geko.NewListFrom([]int{1, 2, 3})

	if !willPanic(func() {
		kol.Set(-1, 0)
	}) {
		t.Fatalf("Set doesn't panic with negative index")
	}

	if !willPanic(func() {
		kol.Set(6, 0)
	}) {
		t.Fatalf("Set doesn't panic with out-of-bound index")
	}

	kol.Set(1, 10)

	if kol.Get(1) != 10 {
		t.Fatalf("Set not effect")
	}
}
