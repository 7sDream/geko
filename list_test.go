package geko_test

import (
	"math/rand"
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

func TestList_Append(t *testing.T) {
	kol := geko.NewList[int]()

	kol.Append(1, 2, 3)

	if !reflect.DeepEqual(kol.List, []int{1, 2, 3}) {
		t.Fatalf("Append not correct")
	}
}

func TestList_Delete(t *testing.T) {
	kol := geko.NewListFrom([]int{1, 2, 3})

	if !willPanic(func() {
		kol.Delete(-1)
	}) {
		t.Fatalf("Delete doesn't panic with negative index")
	}

	if !willPanic(func() {
		kol.Delete(3)
	}) {
		t.Fatalf("Delete doesn't panic with out-of-bound index")
	}

	kol.Delete(1)

	if !reflect.DeepEqual(kol.List, []int{1, 3}) {
		t.Fatalf("Delete not correct")
	}
}

func TestList_Len(t *testing.T) {
	for times := 0; times < 20; times++ {
		kol := geko.NewList[int]()
		exceptedLength := rand.Int() % 100

		for i := 0; i < exceptedLength; i++ {
			kol.Append(i)
		}

		length := kol.Len()
		if length != exceptedLength {
			t.Fatalf("Len excepted %d, got %d", exceptedLength, length)
		}
	}
}
