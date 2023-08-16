package geko_test

import (
	"math/rand"
	"reflect"
	"testing"

	"github.com/7sDream/geko"
)

func TestList_New(t *testing.T) {
	l := geko.NewList[int]()

	if l.List != nil {
		t.Fatalf("NewList inner slice is not nil")
	}

	l2 := geko.NewListFrom([]int{1, 2, 3})

	if !reflect.DeepEqual(l2.List, []int{1, 2, 3}) {
		t.Fatalf("NewListFrom doesn't store origin slice")
	}
}

func TestList_NewWithCapacity(t *testing.T) {
	l := geko.NewListWithCapacity[int](12)

	if cap(l.List) != 12 {
		t.Fatalf("NewListWithCapacity inner slice does not have correct capacity")
	}
}

func TestList_Get(t *testing.T) {
	l := geko.NewListFrom([]int{1, 2, 3})

	if !willPanic(func() {
		l.Get(-1)
	}) {
		t.Fatalf("Get doesn't panic with negative index")
	}

	if !willPanic(func() {
		l.Get(4)
	}) {
		t.Fatalf("Get doesn't panic with out-of-bound index")
	}

	if l.Get(1) != 2 {
		t.Fatalf("Get doesn't return correct value")
	}
}

func TestList_Set(t *testing.T) {
	l := geko.NewListFrom([]int{1, 2, 3})

	if !willPanic(func() {
		l.Set(-1, 0)
	}) {
		t.Fatalf("Set doesn't panic with negative index")
	}

	if !willPanic(func() {
		l.Set(6, 0)
	}) {
		t.Fatalf("Set doesn't panic with out-of-bound index")
	}

	l.Set(1, 10)

	if l.Get(1) != 10 {
		t.Fatalf("Set not effect")
	}
}

func TestList_Append(t *testing.T) {
	l := geko.NewList[int]()

	l.Append(1, 2, 3)

	if !reflect.DeepEqual(l.List, []int{1, 2, 3}) {
		t.Fatalf("Append not correct")
	}
}

func TestList_Delete(t *testing.T) {
	l := geko.NewListFrom([]int{1, 2, 3})

	if !willPanic(func() {
		l.Delete(-1)
	}) {
		t.Fatalf("Delete doesn't panic with negative index")
	}

	if !willPanic(func() {
		l.Delete(3)
	}) {
		t.Fatalf("Delete doesn't panic with out-of-bound index")
	}

	l.Delete(1)

	if !reflect.DeepEqual(l.List, []int{1, 3}) {
		t.Fatalf("Delete not correct")
	}
}

func TestList_Len(t *testing.T) {
	for times := 0; times < 20; times++ {
		l := geko.NewList[int]()
		exceptedLength := rand.Int() % 100

		for i := 0; i < exceptedLength; i++ {
			l.Append(i)
		}

		length := l.Len()
		if length != exceptedLength {
			t.Fatalf("Len excepted %d, got %d", exceptedLength, length)
		}
	}
}
