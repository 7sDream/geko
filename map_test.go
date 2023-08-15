package geko_test

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/7sDream/geko"
)

func TestMap_Get(t *testing.T) {
	kom := geko.NewMap[string, int]()
	kom.Set("one", 1)
	kom.Set("two", 2)

	if v, _ := kom.Get("one"); v != 1 {
		t.Fatalf("Expect %d, got %d", 1, v)
	}

	if v, _ := kom.Get("two"); v != 2 {
		t.Fatalf("Expect %d, got %d", 2, v)
	}

	if _, exist := kom.Get("not_exist"); exist != false {
		t.Fatalf("Get a not exist key should return false")
	}

	if kom.GetOrZeroValue("not_exist") != 0 {
		t.Fatalf("Get a not exist key should return zero value")
	}

	kom2 := geko.NewMap[string, *int]()
	if kom2.GetOrZeroValue("not_exist") != nil {
		t.Fatalf("Get a not exist key should return zero value")
	}
}

func willPanic(f func()) (result bool) {
	defer func() {
		if r := recover(); r != nil {
			result = true
		}
	}()

	f()

	return result
}

func TestMap_GetKeyByIndex(t *testing.T) {
	kom := geko.NewMap[string, int]()

	if !willPanic(func() {
		kom.GetKeyByIndex(0)
	}) {
		t.Fatalf("GetKeyByIndex with empty map didn't panic")
	}

	kom.Set("one", 1)
	kom.Set("three", 2)
	kom.Set("two", 2)
	kom.Set("three", 3)

	if !willPanic(func() {
		kom.GetKeyByIndex(-1)
	}) {
		t.Fatalf("GetKeyByIndex negative index didn't panic")
	}

	if !willPanic(func() {
		kom.GetKeyByIndex(10)
	}) {
		t.Fatalf("GetKeyByIndex out-of-bound index didn't panic")
	}

	expected := "one"
	if v := kom.GetKeyByIndex(0); v != "one" {
		t.Fatalf("GetKeyByIndex(2), Expect %#v, got %#v", expected, v)
	}
}

func TestMap_GetByIndex(t *testing.T) {
	kom := geko.NewMap[string, int]()

	if !willPanic(func() {
		kom.GetByIndex(0)
	}) {
		t.Fatalf("GetByIndex with empty map didn't panic")
	}

	kom.Set("one", 1)
	kom.Set("three", 2)
	kom.Set("two", 2)
	kom.Set("three", 3)

	if !willPanic(func() {
		kom.GetByIndex(-1)
	}) {
		t.Fatalf("GetByIndex negative index didn't panic")
	}

	if !willPanic(func() {
		kom.GetByIndex(10)
	}) {
		t.Fatalf("GetByIndex out-of-bound index didn't panic")
	}

	expected := geko.Pair[string, int]{Key: "three", Value: 3}
	if v := kom.GetByIndex(1); v != expected {
		t.Fatalf("GetByIndex(2), Expect %#v, got %#v", expected, v)
	}
}

func TestMap_GetValueByIndex(t *testing.T) {
	kom := geko.NewMap[string, int]()

	if !willPanic(func() {
		kom.GetValueByIndex(0)
	}) {
		t.Fatalf("GetValueByIndex with empty map didn't panic")
	}

	kom.Set("one", 1)
	kom.Set("three", 2)
	kom.Set("two", 2)
	kom.Set("three", 3)

	if !willPanic(func() {
		kom.GetValueByIndex(-1)
	}) {
		t.Fatalf("GetValueByIndex negative index didn't panic")
	}

	if !willPanic(func() {
		kom.GetValueByIndex(10)
	}) {
		t.Fatalf("GetValueByIndex out-of-bound index didn't panic")
	}

	expected := 2
	if v := kom.GetValueByIndex(2); v != expected {
		t.Fatalf("GetValueByIndex(2), Expect %#v, got %#v", expected, v)
	}
}

func TestMap_Set(t *testing.T) {
	kom := geko.NewMap[string, int]()
	kom.Set("a", 1)
	kom.Set("b", 2)
	kom.Set("b", 3)
	kom.Set("c", 4)
	kom.Set("b", 5)

	keys := []string{
		kom.GetKeyByIndex(0),
		kom.GetKeyByIndex(1),
		kom.GetKeyByIndex(2),
	}
	expectedKeys := []string{"a", "b", "c"}
	if !reflect.DeepEqual(keys, expectedKeys) {
		t.Fatalf("After Set, Expect keys %#v, got %#v", expectedKeys, keys)
	}

	values := []int{
		kom.GetOrZeroValue("a"),
		kom.GetOrZeroValue("b"),
		kom.GetOrZeroValue("c"),
	}
	expectedValues := []int{1, 5, 4}
	if !reflect.DeepEqual(values, expectedValues) {
		t.Fatalf("After Set, Expect keys %#v, got %#v", expectedValues, values)
	}
}

func TestMap_Append(t *testing.T) {
	kom := geko.NewMap[string, int]()
	kom.Append([]geko.Pair[string, int]{
		{"s", 2},
		{"z", 7},
		{"z", 4},
		{"w", 9},
		{"z", 1},
	}...)

	keys := []string{
		kom.GetKeyByIndex(0),
		kom.GetKeyByIndex(1),
		kom.GetKeyByIndex(2),
	}
	expectedKeys := []string{"s", "z", "w"}
	if !reflect.DeepEqual(keys, expectedKeys) {
		t.Fatalf("After Set, Expect keys %#v, got %#v", expectedKeys, keys)
	}

	values := []int{
		kom.GetOrZeroValue("s"),
		kom.GetOrZeroValue("z"),
		kom.GetOrZeroValue("w"),
	}
	expectedValues := []int{2, 1, 9}
	if !reflect.DeepEqual(values, expectedValues) {
		t.Fatalf("After Set, Expect keys %#v, got %#v", expectedValues, values)
	}
}

// Iterate over the map.
func ExampleMap_GetByIndex() {
	kom := geko.NewMap[string, int]()
	kom.Set("one", 1)
	kom.Set("three", 2)
	kom.Set("two", 2)
	kom.Set("three", 3) // do not change order of key "three", it will stay ahead of "two".

	for i := 0; i < kom.Len(); i++ {
		pair := kom.GetByIndex(i)
		fmt.Printf("%s: %d\n", pair.Key, pair.Value)
	}

	data, _ := json.Marshal(kom)
	fmt.Printf("marshal result: %s", string(data))
	// Output:
	// one: 1
	// three: 3
	// two: 2
	// marshal result: {"one":1,"three":3,"two":2}
}
