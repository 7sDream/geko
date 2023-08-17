package geko_test

import (
	"math/rand"
	"reflect"
	"strconv"
	"testing"

	"github.com/7sDream/geko"
)

func TestPairs_New(t *testing.T) {
	l := geko.NewPairs[string, int]()

	if l.List != nil {
		t.Fatalf("NewPairs inner slice is not nil")
	}

	list := []geko.Pair[string, int]{
		{"one", 1},
		{"two", 2},
		{"three", 3},
	}
	l2 := geko.NewPairsFrom(list)

	if !reflect.DeepEqual(l2.List, list) {
		t.Fatalf("NewPairs doesn't store origin slice")
	}
}

func TestPairs_NewWithCapacity(t *testing.T) {
	l := geko.NewPairsWithCapacity[string, int](12)

	if cap(l.List) != 12 {
		t.Fatalf("NewPairsWithCapacity inner slice does not have correct capacity")
	}
}

func TestPairs_Get(t *testing.T) {
	m := geko.NewPairs[string, int]()
	m.Add("one", 1)
	m.Add("two", 2)
	m.Add("two", 22)

	if v := m.Get("one"); v[0] != 1 {
		t.Fatalf("Expect %d, got %d", 1, v)
	}

	value := m.Get("two")
	exceptedValues := []int{2, 22}
	if !reflect.DeepEqual(value, exceptedValues) {
		t.Fatalf("Expect %d, got %d", exceptedValues, value)
	}

	if v := m.Get("not_exist"); len(v) != 0 {
		t.Fatalf("Get a not exist key should return empty slice")
	}
}

func TestPairs_Has(t *testing.T) {
	m := geko.NewPairs[string, int]()
	m.Add("one", 1)
	m.Add("two", 2)

	if !m.Has("one") {
		t.Fatalf("Has said key 'one' does not exist")
	}

	if m.Has("three") {
		t.Fatalf("Has said key 'three' exist")
	}
}

func TestPairs_Count(t *testing.T) {
	m := geko.NewPairs[string, int]()
	m.Add("one", 1)
	m.Add("two", 2)
	m.Add("two", 22)

	if m.Count("zero") != 0 {
		t.Fatalf("Count 'zero' not correct")
	}

	if m.Count("one") != 1 {
		t.Fatalf("Count 'one' not correct")
	}

	if m.Count("two") != 2 {
		t.Fatalf("Count 'two' not correct")
	}
}

func TestPairs_GetXXXOrZeroValue(t *testing.T) {
	m := geko.NewPairs[string, int]()
	m.Add("one", 1)
	m.Add("two", 2)
	m.Add("one", 11)
	m.Add("two", 22)

	if v := m.GetFirstOrZeroValue("one"); v != 1 {
		t.Fatalf("Expect %d, got %d", 1, v)
	}

	if v := m.GetLastOrZeroValue("two"); v != 22 {
		t.Fatalf("Expect %d, got %d", 22, v)
	}

	if v := m.GetFirstOrZeroValue("not_exist"); v != 0 {
		t.Fatalf("GetFirstOrZeroValue a not exist key should return zero value")
	}

	if v := m.GetLastOrZeroValue("not_exist"); v != 0 {
		t.Fatalf("GetLastOrZeroValue a not exist key should return zero value")
	}
}

func TestPairs_GetKeyByIndex(t *testing.T) {
	m := geko.NewPairs[string, int]()

	if !willPanic(func() {
		m.GetKeyByIndex(0)
	}) {
		t.Fatalf("GetKeyByIndex with empty map didn't panic")
	}

	m.Add("one", 1)
	m.Add("three", 2)
	m.Add("two", 2)
	m.Add("three", 3)

	if !willPanic(func() {
		m.GetKeyByIndex(-1)
	}) {
		t.Fatalf("GetKeyByIndex with negative index didn't panic")
	}

	if !willPanic(func() {
		m.GetKeyByIndex(10)
	}) {
		t.Fatalf("GetKeyByIndex with out-of-bound index didn't panic")
	}

	expected := "three"
	if v := m.GetKeyByIndex(3); v != expected {
		t.Fatalf("GetKeyByIndex(3), Expect %#v, got %#v", expected, v)
	}
}

func TestPairs_GetByIndex(t *testing.T) {
	m := geko.NewPairs[string, int]()

	if !willPanic(func() {
		m.GetByIndex(0)
	}) {
		t.Fatalf("GetByIndex with empty map didn't panic")
	}

	m.Add("one", 1)
	m.Add("three", 2)
	m.Add("two", 2)
	m.Add("three", 3)

	if !willPanic(func() {
		m.GetByIndex(-1)
	}) {
		t.Fatalf("GetByIndex negative index didn't panic")
	}

	if !willPanic(func() {
		m.GetByIndex(10)
	}) {
		t.Fatalf("GetByIndex out-of-bound index didn't panic")
	}

	expected := geko.Pair[string, int]{Key: "three", Value: 2}
	if v := m.GetByIndex(1); v != expected {
		t.Fatalf("GetByIndex(1), Expect %#v, got %#v", expected, v)
	}
}

func TestPairs_SetKeyByIndex(t *testing.T) {
	m := geko.NewPairs[string, int]()

	if !willPanic(func() {
		m.SetKeyByIndex(0, "new")
	}) {
		t.Fatalf("SetKeyByIndex with empty map didn't panic")
	}

	m.Add("one", 1)
	m.Add("four", 2)
	m.Add("three", 3)

	if !willPanic(func() {
		m.SetKeyByIndex(3, "new")
	}) {
		t.Fatalf("SetKeyByIndex with out-of-bound index didn't panic")
	}

	m.SetKeyByIndex(1, "two")
	if m.GetKeyByIndex(1) != "two" {
		t.Fatalf("SetKeyByIndex do not effect")
	}
}

func TestPairs_SetValueByIndex(t *testing.T) {
	m := geko.NewPairs[string, int]()

	if !willPanic(func() {
		m.SetValueByIndex(0, 0)
	}) {
		t.Fatalf("SetValueByIndex with empty map didn't panic")
	}

	m.Add("one", 1)
	m.Add("two", 4)
	m.Add("three", 3)

	if !willPanic(func() {
		m.SetKeyByIndex(3, "new")
	}) {
		t.Fatalf("SetValueByIndex with out-of-bound index didn't panic")
	}

	m.SetValueByIndex(1, 2)
	if m.GetValueByIndex(1) != 2 {
		t.Fatalf("SetValueByIndex do not effect")
	}
}

func TestPairs_SetByIndex(t *testing.T) {
	m := geko.NewPairs[string, int]()

	if !willPanic(func() {
		m.SetByIndex(0, "zero", 0)
	}) {
		t.Fatalf("SetByIndex with empty map didn't panic")
	}

	m.Add("one", 1)
	m.Add("four", 4)
	m.Add("three", 3)

	if !willPanic(func() {
		m.SetByIndex(3, "new", 0)
	}) {
		t.Fatalf("SetByIndex with out-of-bound index didn't panic")
	}

	m.SetByIndex(1, "two", 2)
	if m.GetByIndex(1) != geko.CreatePair("two", 2) {
		t.Fatalf("SetByIndex do not effect")
	}
}

func TestPairs_GetValueByIndex(t *testing.T) {
	m := geko.NewPairs[string, int]()

	if !willPanic(func() {
		m.GetValueByIndex(0)
	}) {
		t.Fatalf("GetValueByIndex with empty map didn't panic")
	}

	m.Add("one", 1)
	m.Add("three", 2)
	m.Add("two", 2)
	m.Add("three", 3)

	if !willPanic(func() {
		m.GetValueByIndex(-1)
	}) {
		t.Fatalf("GetValueByIndex negative index didn't panic")
	}

	if !willPanic(func() {
		m.GetValueByIndex(10)
	}) {
		t.Fatalf("GetValueByIndex out-of-bound index didn't panic")
	}

	expected := 2
	if v := m.GetValueByIndex(2); v != expected {
		t.Fatalf("GetValueByIndex(2), Expect %#v, got %#v", expected, v)
	}
}

func TestPairs_Add(t *testing.T) {
	m := geko.NewPairs[string, int]()
	m.Add("a", 1)
	m.Add("b", 2)
	m.Add("a", 3)

	keys := m.Keys()
	values := m.Values()

	exceptedKeys := []string{"a", "b", "a"}
	if !reflect.DeepEqual(keys, exceptedKeys) {
		t.Fatalf(
			"Excepted keys %#v, got %#v",
			exceptedKeys, keys,
		)
	}

	exceptedValues := []int{1, 2, 3}
	if !reflect.DeepEqual(values, exceptedValues) {
		t.Fatalf(
			"Excepted values %#v, got %#v",
			exceptedValues, values,
		)
	}
}

func TestPairs_Append(t *testing.T) {
	m := geko.NewPairs[string, int]()
	m.Append([]geko.Pair[string, int]{
		{"s", 2},
		{"z", 7},
		{"z", 4},
		{"w", 9},
		{"z", 1},
	}...)

	keys := m.Keys()
	expectedKeys := []string{"s", "z", "z", "w", "z"}
	if !reflect.DeepEqual(keys, expectedKeys) {
		t.Fatalf("After Append, expect keys %#v, got %#v", expectedKeys, keys)
	}

	values := [][]int{
		m.Get("s"),
		m.Get("z"),
		m.Get("w"),
	}
	expectedValues := [][]int{{2}, {7, 4, 1}, {9}}
	if !reflect.DeepEqual(values, expectedValues) {
		t.Fatalf("After Append, expect keys %#v, got %#v", expectedValues, values)
	}
}

func TestPairs_Delete(t *testing.T) {
	m := geko.NewPairs[string, int]()
	m.Add("a", 1)

	m.Delete("b") // should not panic

	m.Delete("a")

	if m.Len() != 0 {
		t.Fatalf("After Delete all item, Map is not empty")
	}

	m = geko.NewPairs[string, int]()
	m.Add("a", 1)
	m.Add("b", 2)
	m.Add("c", 3)
	m.Add("b", 4)
	m.Delete("b")

	if m.Len() != 2 {
		t.Fatalf("After Delete item, Len does not correct")
	}

	if m.Count("b") != 0 {
		t.Fatalf("Delete do not delete all matched item")
	}
}

func TestPairs_DeleteByIndex(t *testing.T) {
	m := geko.NewPairs[string, int]()

	if !willPanic(func() {
		m.DeleteByIndex(1)
	}) {
		t.Fatalf("DeleteByIndex with empty map didn't panic")
	}

	m.Add("a", 1)
	m.Add("b", 2)
	m.Add("b", 22)
	m.Add("c", 3)

	m.DeleteByIndex(1)

	if m.Len() != 3 {
		t.Fatalf("After DeleteByIndex, Len does not correct")
	}

	if !m.Has("b") {
		t.Fatalf("After DeleteByIndex, all same key deleted")
	}

	m.DeleteByIndex(1)

	keys := m.Keys()
	excepted := []string{"a", "c"}
	if !reflect.DeepEqual(keys, excepted) {
		t.Fatalf("After another DeleteByIndex, excepted keys %#v, got %#v", excepted, keys)
	}
}

func TestPairs_Clear(t *testing.T) {
	m := geko.NewPairs[string, int]()
	m.Add("a", 1)
	m.Add("b", 2)
	m.Clear()

	if m.Len() != 0 {
		t.Fatalf("After Clean, map is not empty")
	}

	// After Clear, new Add should not panic
	m.Add("b", 2)
	m.Add("a", 1)
	keys := m.Keys()
	excepted := []string{"b", "a"}
	if !reflect.DeepEqual(keys, excepted) {
		t.Fatalf("After Clean, old values should not effect new order")
	}
}

func TestPairs_Len(t *testing.T) {
	for times := 0; times < 20; times++ {
		exceptedLength := rand.Int() % 100

		m := geko.NewPairs[string, int]()
		for i := 0; i < exceptedLength; i++ {
			m.Add(strconv.Itoa(i), i)
		}

		length := m.Len()
		if length != exceptedLength {
			t.Fatalf("Length excepted %d, got %d", exceptedLength, length)
		}
	}
}

func TestPairs_Keys(t *testing.T) {
	m := geko.NewPairs[string, int]()
	m.Add("one", 1)
	m.Add("three", 2)
	m.Add("two", 2)
	m.Add("three", 3)

	m.Delete("one")

	excepted := []string{"three", "two", "three"}
	keys := m.Keys()
	if !reflect.DeepEqual(keys, excepted) {
		t.Fatalf("Excepted keys %#v, got %#v", excepted, keys)
	}

	keys[0] = "haha"
	if reflect.DeepEqual(keys, m.Keys()) {
		t.Fatalf("Modify return keys should not effect pairs")
	}
}

func TestPairs_Values(t *testing.T) {
	m := geko.NewPairs[string, int]()
	m.Add("one", 1)
	m.Add("three", 2)
	m.Add("two", 2)
	m.Add("three", 3)

	m.Delete("one")

	excepted := []int{2, 2, 3}
	values := m.Values()
	if !reflect.DeepEqual(values, excepted) {
		t.Fatalf("Excepted values %#v, got %#v", excepted, values)
	}

	values[0] = 100
	if reflect.DeepEqual(values, m.Values()) {
		t.Fatalf("Modify return values should not effect map")
	}

	type s struct {
		Value int
	}

	m2 := geko.NewPairs[string, *s]()
	m2.Add("one", &s{Value: 1})
	m2.Add("two", &s{Value: 2})
	m2.Add("three", &s{Value: 3})

	m2.Values()[2].Value = 100

	if m2.GetFirstOrZeroValue("three").Value != 100 {
		t.Fatalf("Use pointer as value type will allow user modifier inner value")
	}
}

func TestPairs_ToMap(t *testing.T) {
	ps := geko.NewPairs[string, int]()
	ps.Add("one", 1)
	ps.Add("three", 2)
	ps.Add("two", 2)
	ps.Add("three", 3)

	cases := []struct {
		strategy geko.DuplicatedKeyStrategy
		keys     []string
		values   []int
	}{
		{geko.UpdateValueKeepOrder, []string{"one", "three", "two"}, []int{1, 3, 2}},
		{geko.UpdateValueUpdateOrder, []string{"one", "two", "three"}, []int{1, 2, 3}},
		{geko.KeepValueUpdateOrder, []string{"one", "two", "three"}, []int{1, 2, 2}},
		{geko.Ignore, []string{"one", "three", "two"}, []int{1, 2, 2}},
	}

	for _, tt := range cases {
		m := ps.ToMap(tt.strategy)

		keys := m.Keys()
		if !reflect.DeepEqual(keys, tt.keys) {
			t.Fatalf(
				"Strategy %d, excepted keys %#v, got %#v",
				tt.strategy, tt.keys, keys,
			)
		}

		values := m.Values()
		if !reflect.DeepEqual(values, tt.values) {
			t.Fatalf(
				"Strategy %d, excepted values %#v, got %#v",
				tt.strategy, tt.values, values,
			)
		}
	}
}

func TestPairs_Dedup(t *testing.T) {
	cases := []struct {
		strategy geko.DedupStrategy
		keys     []string
		values   []int
	}{
		{geko.KeepFirst, []string{"one", "three", "two"}, []int{1, 2, 2}},
		{geko.KeepLast, []string{"one", "two", "three"}, []int{1, 2, 3}},

		/* invalid value treat as default strategy */
		{geko.DedupStrategy(10), []string{"one", "three", "two"}, []int{1, 2, 2}},
	}

	for _, tt := range cases {
		ps := geko.NewPairs[string, int]()
		ps.Add("one", 1)
		ps.Add("three", 2)
		ps.Add("two", 2)
		ps.Add("three", 3)
		ps.Dedup(tt.strategy)

		keys := ps.Keys()
		if !reflect.DeepEqual(keys, tt.keys) {
			t.Fatalf(
				"Strategy %d, excepted keys %#v, got %#v",
				tt.strategy, tt.keys, keys,
			)
		}

		values := ps.Values()
		if !reflect.DeepEqual(values, tt.values) {
			t.Fatalf(
				"Strategy %d, excepted values %#v, got %#v",
				tt.strategy, tt.values, values,
			)
		}
	}
}

func TestPairs_Sort(t *testing.T) {
	m := geko.NewPairs[int, string]()
	m.Add(3, "three.2")
	m.Add(1, "one")
	m.Add(4, "four")
	m.Add(2, "two")
	m.Add(3, "three.1")

	m.Sort(func(a, b *geko.Pair[int, string]) bool {
		return a.Key < b.Key
	})

	exceptedPairs := []geko.Pair[int, string]{
		{1, "one"},
		{2, "two"},
		{3, "three.2"},
		{3, "three.1"}, // Test sort stability
		{4, "four"},
	}

	if !reflect.DeepEqual(m.List, exceptedPairs) {
		t.Fatalf("Sort result excepted %#v, got %#v", exceptedPairs, m.List)
	}
}

func TestPairs_Filter(t *testing.T) {
	m := geko.NewPairs[int, string]()
	m.Add(1, "one")
	m.Add(2, "two")
	m.Add(3, "three")
	m.Add(4, "four")

	m.Filter(func(p *geko.Pair[int, string]) bool {
		return p.Key%2 == 0
	})

	exceptedPairs := []geko.Pair[int, string]{
		{2, "two"},
		{4, "four"},
	}

	if !reflect.DeepEqual(exceptedPairs, m.List) {
		t.Fatalf("Filter result excepted %#v, got %#v", exceptedPairs, m.List)
	}
}
