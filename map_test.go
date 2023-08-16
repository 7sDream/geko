package geko_test

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"reflect"
	"strconv"
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
	for _, strategy := range []geko.DuplicateKeyStrategy{
		geko.UpdateValueKeepOrder,
		geko.UpdateValueUpdateOrder,
		geko.KeepValueUpdateOrder,
		geko.Ignore,
	} {
		kom := geko.NewMap[string, int]()
		kom.SetDuplicateKeyStrategy(strategy)
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
}

func TestMap_Add(t *testing.T) {
	cases := []struct {
		strategy       geko.DuplicateKeyStrategy
		exceptedKeys   []string
		exceptedValues []int
	}{
		{geko.UpdateValueKeepOrder, []string{"a", "b"}, []int{3, 2}},
		{geko.UpdateValueUpdateOrder, []string{"b", "a"}, []int{2, 3}},
		{geko.KeepValueUpdateOrder, []string{"b", "a"}, []int{2, 1}},
		{geko.Ignore, []string{"a", "b"}, []int{1, 2}},
	}

	for _, tt := range cases {
		kom := geko.NewMap[string, int]()
		kom.SetDuplicateKeyStrategy(tt.strategy)
		kom.Add("a", 1)
		kom.Add("b", 2)
		kom.Add("a", 3)

		if strategy := kom.DuplicateKeyStrategy(); strategy != tt.strategy {
			t.Fatalf(
				"Excepted strategy %#v, got %#v",
				tt.strategy, strategy,
			)
		}

		keys := kom.Keys()
		values := kom.Values()

		if !reflect.DeepEqual(keys, tt.exceptedKeys) {
			t.Fatalf(
				"for strategy %#v, Excepted keys %#v, got %#v",
				tt.strategy, tt.exceptedKeys, keys,
			)
		}

		if !reflect.DeepEqual(values, tt.exceptedValues) {
			t.Fatalf(
				"for strategy %#v, Excepted values %#v, got %#v",
				tt.strategy, tt.exceptedValues, values,
			)
		}
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

func TestMap_Delete(t *testing.T) {
	kom := geko.NewMap[string, int]()
	kom.Set("a", 1)

	kom.Delete("b") // should not panic

	kom.Delete("a")

	if kom.Len() != 0 {
		t.Fatalf("After Delete all item, Map is not empty")
	}

	kom = geko.NewMap[string, int]()
	kom.Set("a", 1)
	kom.Set("b", 2)
	kom.Set("c", 3)
	kom.Delete("b")

	if kom.Len() != 2 {
		t.Fatalf("After Delete a item, Len does not correct")
	}

	if _, exist := kom.Get("b"); exist != false {
		t.Fatalf("After Delete item, it still exist")
	}

	kom.Set("b", 4)

	if kom.Len() != 3 {
		t.Fatalf("After Delete and Set a same key, Len does not correct")
	}

	if v := kom.GetValueByIndex(2); v != 4 {
		t.Fatalf("Item does not appear in last after Delete and Set")
	}
}

func TestMap_DeleteByIndex(t *testing.T) {
	kom := geko.NewMap[string, int]()

	if !willPanic(func() {
		kom.DeleteByIndex(1)
	}) {
		t.Fatalf("DeleteByIndex with empty map didn't panic")
	}

	kom.Set("a", 1)
	kom.Set("b", 2)
	kom.Set("c", 3)
	kom.DeleteByIndex(1)

	if kom.Len() != 2 {
		t.Fatalf("After DeleteByIndex, Len does not correct")
	}

	if _, exist := kom.Get("b"); exist {
		t.Fatalf("After DeleteByIndex, it still exist")
	}

	keys := kom.Keys()
	excepted := []string{"a", "c"}
	if !reflect.DeepEqual(keys, excepted) {
		t.Fatalf("After DeleteByIndex, excepted keys %#v, got %#v", excepted, keys)
	}
}

func TestMap_Clear(t *testing.T) {
	kom := geko.NewMap[string, int]()
	kom.Set("a", 1)
	kom.Set("b", 2)
	kom.Clear()

	if kom.Len() != 0 {
		t.Fatalf("After Clean, map is not empty")
	}

	if len(kom.Keys()) != 0 {
		t.Fatalf("After Clean, map Keys() is not empty")
	}

	// After Clear, new Set should not panic
	kom.Set("b", 2)
	kom.Set("a", 1)
	keys := kom.Keys()
	excepted := []string{"b", "a"}
	if !reflect.DeepEqual(keys, excepted) {
		t.Fatalf("After Clean, old values should not effect new order")
	}
}

func TestMap_Len(t *testing.T) {
	for times := 0; times < 20; times++ {
		exceptedLength := rand.Int() % 100

		kom := geko.NewMap[string, int]()
		for i := 0; i < exceptedLength; i++ {
			kom.Set(strconv.Itoa(i), i)
		}

		length := kom.Len()
		if length != exceptedLength {
			t.Fatalf("Length excepted %d, got %d", exceptedLength, length)
		}
	}
}

func TestMap_Keys(t *testing.T) {
	kom := geko.NewMap[string, int]()
	kom.Set("one", 1)
	kom.Set("three", 2)
	kom.Set("two", 2)
	kom.Set("three", 3)

	kom.Delete("one")

	excepted := []string{"three", "two"}
	keys := kom.Keys()
	if !reflect.DeepEqual(keys, excepted) {
		t.Fatalf("Excepted keys %#v, got %#v", excepted, keys)
	}

	keys[0] = "haha"
	if reflect.DeepEqual(keys, kom.Keys()) {
		t.Fatalf("Modify return keys should not effect map")
	}
}

func TestMap_Values(t *testing.T) {
	kom := geko.NewMap[string, int]()
	kom.Set("one", 1)
	kom.Set("three", 2)
	kom.Set("two", 2)
	kom.Set("three", 3)

	kom.Delete("one")

	excepted := []int{3, 2}
	values := kom.Values()
	if !reflect.DeepEqual(values, excepted) {
		t.Fatalf("Excepted values %#v, got %#v", excepted, values)
	}

	values[0] = 100
	if reflect.DeepEqual(values, kom.Values()) {
		t.Fatalf("Modify return values should not effect map")
	}

	type s struct {
		Value int
	}

	kom2 := geko.NewMap[string, *s]()
	kom2.Set("one", &s{Value: 1})
	kom2.Set("two", &s{Value: 2})
	kom2.Set("three", &s{Value: 3})

	kom2.Values()[2].Value = 100

	if kom2.GetOrZeroValue("three").Value != 100 {
		t.Fatalf("Use pointer as value type will allow user modifier inner value")
	}
}

func TestMap_Pairs(t *testing.T) {
	kom := geko.NewMap[string, int]()
	kom.Set("one", 1)
	kom.Set("three", 2)
	kom.Set("two", 2)
	kom.Set("three", 3)
	kom.Delete("one")

	expected := []geko.Pair[string, int]{
		{"three", 3},
		{"two", 2},
	}
	pairs := kom.Pairs().List
	if !reflect.DeepEqual(pairs, expected) {
		t.Fatalf("Excepted %#v, got %#v", expected, pairs)
	}
}

func TestMap_Sort(t *testing.T) {
	kom := geko.NewMap[int, string]()
	kom.Set(3, "three")
	kom.Set(1, "one")
	kom.Set(4, "four")
	kom.Set(2, "two")

	kom.Sort(func(a, b *geko.Pair[int, string]) bool {
		return a.Key < b.Key
	})

	exceptedPairs := []geko.Pair[int, string]{
		{1, "one"},
		{2, "two"},
		{3, "three"},
		{4, "four"},
	}

	pairs := kom.Pairs().List

	if !reflect.DeepEqual(pairs, exceptedPairs) {
		t.Fatalf("Sort result excepted %#v, got %#v", exceptedPairs, pairs)
	}
}

func TestMap_Filter(t *testing.T) {
	kom := geko.NewMap[int, string]()
	kom.Set(1, "one")
	kom.Set(2, "two")
	kom.Set(3, "three")
	kom.Set(4, "four")

	kom.Filter(func(p *geko.Pair[int, string]) bool {
		return p.Key%2 == 0
	})

	exceptedPairs := []geko.Pair[int, string]{
		{2, "two"},
		{4, "four"},
	}

	pairs := kom.Pairs().List

	if !reflect.DeepEqual(pairs, exceptedPairs) {
		t.Fatalf("Filter result excepted %#v, got %#v", exceptedPairs, pairs)
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
