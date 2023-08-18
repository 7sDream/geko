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

func ExamplePairs() {
	m := geko.NewPairs[string, int]()

	m.Add("one", 1)
	m.Add("three", 2)
	m.Add("two", 2)
	m.Add("three", 3)
	for i, length := 0, m.Len(); i < length; i++ {
		pair := m.GetByIndex(i)
		fmt.Printf("%s: %d\n", pair.Key, pair.Value)
	}

	fmt.Println("-----")

	m.Dedup(geko.Ignore)
	for i, length := 0, m.Len(); i < length; i++ {
		pair := m.GetByIndex(i)
		fmt.Printf("%s: %d\n", pair.Key, pair.Value)
	}

	// Output:
	// one: 1
	// three: 2
	// two: 2
	// three: 3
	// -----
	// one: 1
	// three: 2
	// two: 2
}

func TestPairs_New(t *testing.T) {
	ps := geko.NewPairs[string, int]()

	if ps.List != nil {
		t.Fatalf("NewPairs inner slice is not nil")
	}

	list := []geko.Pair[string, int]{
		{"one", 1},
		{"two", 2},
		{"three", 3},
	}
	ps2 := geko.NewPairsFrom(list)

	if !reflect.DeepEqual(ps2.List, list) {
		t.Fatalf("NewPairs doesn't store origin slice")
	}
}

func TestPairs_NewWithCapacity(t *testing.T) {
	ps := geko.NewPairsWithCapacity[string, int](12)

	if cap(ps.List) != 12 {
		t.Fatalf("NewPairsWithCapacity inner slice does not have correct capacity")
	}
}

func TestPairs_Get(t *testing.T) {
	ps := geko.NewPairs[string, int]()
	ps.Add("one", 1)
	ps.Add("two", 2)
	ps.Add("two", 22)

	if v := ps.Get("one"); v[0] != 1 {
		t.Fatalf("Expect %d, got %d", 1, v)
	}

	value := ps.Get("two")
	exceptedValues := []int{2, 22}
	if !reflect.DeepEqual(value, exceptedValues) {
		t.Fatalf("Expect %d, got %d", exceptedValues, value)
	}

	if v := ps.Get("not_exist"); len(v) != 0 {
		t.Fatalf("Get a not exist key should return empty slice")
	}
}

func TestPairs_Has(t *testing.T) {
	ps := geko.NewPairs[string, int]()
	ps.Add("one", 1)
	ps.Add("two", 2)

	if !ps.Has("one") {
		t.Fatalf("Has said key 'one' does not exist")
	}

	if ps.Has("three") {
		t.Fatalf("Has said key 'three' exist")
	}
}

func TestPairs_Count(t *testing.T) {
	ps := geko.NewPairs[string, int]()
	ps.Add("one", 1)
	ps.Add("two", 2)
	ps.Add("two", 22)

	if ps.Count("zero") != 0 {
		t.Fatalf("Count 'zero' not correct")
	}

	if ps.Count("one") != 1 {
		t.Fatalf("Count 'one' not correct")
	}

	if ps.Count("two") != 2 {
		t.Fatalf("Count 'two' not correct")
	}
}

func TestPairs_GetXXXOrZeroValue(t *testing.T) {
	ps := geko.NewPairs[string, int]()
	ps.Add("one", 1)
	ps.Add("two", 2)
	ps.Add("one", 11)
	ps.Add("two", 22)

	if v := ps.GetFirstOrZeroValue("one"); v != 1 {
		t.Fatalf("Expect %d, got %d", 1, v)
	}

	if v := ps.GetLastOrZeroValue("two"); v != 22 {
		t.Fatalf("Expect %d, got %d", 22, v)
	}

	if v := ps.GetFirstOrZeroValue("not_exist"); v != 0 {
		t.Fatalf("GetFirstOrZeroValue a not exist key should return zero value")
	}

	if v := ps.GetLastOrZeroValue("not_exist"); v != 0 {
		t.Fatalf("GetLastOrZeroValue a not exist key should return zero value")
	}
}

func TestPairs_GetKeyByIndex(t *testing.T) {
	ps := geko.NewPairs[string, int]()

	if !willPanic(func() {
		ps.GetKeyByIndex(0)
	}) {
		t.Fatalf("GetKeyByIndex with empty map didn't panic")
	}

	ps.Add("one", 1)
	ps.Add("three", 2)
	ps.Add("two", 2)
	ps.Add("three", 3)

	if !willPanic(func() {
		ps.GetKeyByIndex(-1)
	}) {
		t.Fatalf("GetKeyByIndex with negative index didn't panic")
	}

	if !willPanic(func() {
		ps.GetKeyByIndex(10)
	}) {
		t.Fatalf("GetKeyByIndex with out-of-bound index didn't panic")
	}

	expected := "three"
	if v := ps.GetKeyByIndex(3); v != expected {
		t.Fatalf("GetKeyByIndex(3), Expect %#v, got %#v", expected, v)
	}
}

func TestPairs_GetByIndex(t *testing.T) {
	ps := geko.NewPairs[string, int]()

	if !willPanic(func() {
		ps.GetByIndex(0)
	}) {
		t.Fatalf("GetByIndex with empty map didn't panic")
	}

	ps.Add("one", 1)
	ps.Add("three", 2)
	ps.Add("two", 2)
	ps.Add("three", 3)

	if !willPanic(func() {
		ps.GetByIndex(-1)
	}) {
		t.Fatalf("GetByIndex negative index didn't panic")
	}

	if !willPanic(func() {
		ps.GetByIndex(10)
	}) {
		t.Fatalf("GetByIndex out-of-bound index didn't panic")
	}

	expected := geko.Pair[string, int]{Key: "three", Value: 2}
	if v := ps.GetByIndex(1); v != expected {
		t.Fatalf("GetByIndex(1), Expect %#v, got %#v", expected, v)
	}
}

func TestPairs_SetKeyByIndex(t *testing.T) {
	ps := geko.NewPairs[string, int]()

	if !willPanic(func() {
		ps.SetKeyByIndex(0, "new")
	}) {
		t.Fatalf("SetKeyByIndex with empty map didn't panic")
	}

	ps.Add("one", 1)
	ps.Add("four", 2)
	ps.Add("three", 3)

	if !willPanic(func() {
		ps.SetKeyByIndex(3, "new")
	}) {
		t.Fatalf("SetKeyByIndex with out-of-bound index didn't panic")
	}

	ps.SetKeyByIndex(1, "two")
	if ps.GetKeyByIndex(1) != "two" {
		t.Fatalf("SetKeyByIndex do not effect")
	}
}

func TestPairs_SetValueByIndex(t *testing.T) {
	ps := geko.NewPairs[string, int]()

	if !willPanic(func() {
		ps.SetValueByIndex(0, 0)
	}) {
		t.Fatalf("SetValueByIndex with empty map didn't panic")
	}

	ps.Add("one", 1)
	ps.Add("two", 4)
	ps.Add("three", 3)

	if !willPanic(func() {
		ps.SetKeyByIndex(3, "new")
	}) {
		t.Fatalf("SetValueByIndex with out-of-bound index didn't panic")
	}

	ps.SetValueByIndex(1, 2)
	if ps.GetValueByIndex(1) != 2 {
		t.Fatalf("SetValueByIndex do not effect")
	}
}

func TestPairs_SetByIndex(t *testing.T) {
	ps := geko.NewPairs[string, int]()

	if !willPanic(func() {
		ps.SetByIndex(0, "zero", 0)
	}) {
		t.Fatalf("SetByIndex with empty map didn't panic")
	}

	ps.Add("one", 1)
	ps.Add("four", 4)
	ps.Add("three", 3)

	if !willPanic(func() {
		ps.SetByIndex(3, "new", 0)
	}) {
		t.Fatalf("SetByIndex with out-of-bound index didn't panic")
	}

	ps.SetByIndex(1, "two", 2)
	if ps.GetByIndex(1) != geko.CreatePair("two", 2) {
		t.Fatalf("SetByIndex do not effect")
	}
}

func TestPairs_GetValueByIndex(t *testing.T) {
	ps := geko.NewPairs[string, int]()

	if !willPanic(func() {
		ps.GetValueByIndex(0)
	}) {
		t.Fatalf("GetValueByIndex with empty map didn't panic")
	}

	ps.Add("one", 1)
	ps.Add("three", 2)
	ps.Add("two", 2)
	ps.Add("three", 3)

	if !willPanic(func() {
		ps.GetValueByIndex(-1)
	}) {
		t.Fatalf("GetValueByIndex negative index didn't panic")
	}

	if !willPanic(func() {
		ps.GetValueByIndex(10)
	}) {
		t.Fatalf("GetValueByIndex out-of-bound index didn't panic")
	}

	expected := 2
	if v := ps.GetValueByIndex(2); v != expected {
		t.Fatalf("GetValueByIndex(2), Expect %#v, got %#v", expected, v)
	}
}

func TestPairs_Add(t *testing.T) {
	ps := geko.NewPairs[string, int]()
	ps.Add("a", 1)
	ps.Add("b", 2)
	ps.Add("a", 3)

	keys := ps.Keys()
	values := ps.Values()

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
	ps := geko.NewPairs[string, int]()
	ps.Append([]geko.Pair[string, int]{
		{"s", 2},
		{"z", 7},
		{"z", 4},
		{"w", 9},
		{"z", 1},
	}...)

	keys := ps.Keys()
	expectedKeys := []string{"s", "z", "z", "w", "z"}
	if !reflect.DeepEqual(keys, expectedKeys) {
		t.Fatalf("After Append, expect keys %#v, got %#v", expectedKeys, keys)
	}

	values := [][]int{
		ps.Get("s"),
		ps.Get("z"),
		ps.Get("w"),
	}
	expectedValues := [][]int{{2}, {7, 4, 1}, {9}}
	if !reflect.DeepEqual(values, expectedValues) {
		t.Fatalf("After Append, expect keys %#v, got %#v", expectedValues, values)
	}
}

func TestPairs_Delete(t *testing.T) {
	ps := geko.NewPairs[string, int]()
	ps.Add("a", 1)

	ps.Delete("b") // should not panic

	ps.Delete("a")

	if ps.Len() != 0 {
		t.Fatalf("After Delete all item, Map is not empty")
	}

	ps = geko.NewPairs[string, int]()
	ps.Add("a", 1)
	ps.Add("b", 2)
	ps.Add("c", 3)
	ps.Add("b", 4)
	ps.Delete("b")

	if ps.Len() != 2 {
		t.Fatalf("After Delete item, Len does not correct")
	}

	if ps.Count("b") != 0 {
		t.Fatalf("Delete do not delete all matched item")
	}
}

func TestPairs_DeleteByIndex(t *testing.T) {
	ps := geko.NewPairs[string, int]()

	if !willPanic(func() {
		ps.DeleteByIndex(1)
	}) {
		t.Fatalf("DeleteByIndex with empty map didn't panic")
	}

	ps.Add("a", 1)
	ps.Add("b", 2)
	ps.Add("b", 22)
	ps.Add("c", 3)

	ps.DeleteByIndex(1)

	if ps.Len() != 3 {
		t.Fatalf("After DeleteByIndex, Len does not correct")
	}

	if !ps.Has("b") {
		t.Fatalf("After DeleteByIndex, all same key deleted")
	}

	ps.DeleteByIndex(1)

	keys := ps.Keys()
	excepted := []string{"a", "c"}
	if !reflect.DeepEqual(keys, excepted) {
		t.Fatalf("After another DeleteByIndex, excepted keys %#v, got %#v", excepted, keys)
	}
}

func TestPairs_Clear(t *testing.T) {
	ps := geko.NewPairs[string, int]()
	ps.Add("a", 1)
	ps.Add("b", 2)
	ps.Clear()

	if ps.Len() != 0 {
		t.Fatalf("After Clean, map is not empty")
	}

	// After Clear, new Add should not panic
	ps.Add("b", 2)
	ps.Add("a", 1)
	keys := ps.Keys()
	excepted := []string{"b", "a"}
	if !reflect.DeepEqual(keys, excepted) {
		t.Fatalf("After Clean, old values should not effect new order")
	}
}

func TestPairs_Len(t *testing.T) {
	for times := 0; times < 20; times++ {
		exceptedLength := rand.Int() % 100

		ps := geko.NewPairs[string, int]()
		for i := 0; i < exceptedLength; i++ {
			ps.Add(strconv.Itoa(i), i)
		}

		length := ps.Len()
		if length != exceptedLength {
			t.Fatalf("Length excepted %d, got %d", exceptedLength, length)
		}
	}
}

func TestPairs_Keys(t *testing.T) {
	ps := geko.NewPairs[string, int]()
	ps.Add("one", 1)
	ps.Add("three", 2)
	ps.Add("two", 2)
	ps.Add("three", 3)

	ps.Delete("one")

	excepted := []string{"three", "two", "three"}
	keys := ps.Keys()
	if !reflect.DeepEqual(keys, excepted) {
		t.Fatalf("Excepted keys %#v, got %#v", excepted, keys)
	}

	keys[0] = "haha"
	if reflect.DeepEqual(keys, ps.Keys()) {
		t.Fatalf("Modify return keys should not effect pairs")
	}
}

func TestPairs_Values(t *testing.T) {
	ps := geko.NewPairs[string, int]()
	ps.Add("one", 1)
	ps.Add("three", 2)
	ps.Add("two", 2)
	ps.Add("three", 3)

	ps.Delete("one")

	excepted := []int{2, 2, 3}
	values := ps.Values()
	if !reflect.DeepEqual(values, excepted) {
		t.Fatalf("Excepted values %#v, got %#v", excepted, values)
	}

	values[0] = 100
	if reflect.DeepEqual(values, ps.Values()) {
		t.Fatalf("Modify return values should not effect map")
	}

	type s struct {
		Value int
	}

	ps2 := geko.NewPairs[string, *s]()
	ps2.Add("one", &s{Value: 1})
	ps2.Add("two", &s{Value: 2})
	ps2.Add("three", &s{Value: 3})

	ps2.Values()[2].Value = 100

	if ps2.GetFirstOrZeroValue("three").Value != 100 {
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
		strategy geko.DuplicatedKeyStrategy
		keys     []string
		values   []int
	}{
		{geko.UpdateValueKeepOrder, []string{"one", "three", "two"}, []int{1, 3, 2}},
		{geko.UpdateValueUpdateOrder, []string{"one", "two", "three"}, []int{1, 2, 3}},
		{geko.KeepValueUpdateOrder, []string{"one", "two", "three"}, []int{1, 2, 2}},
		{geko.Ignore, []string{"one", "three", "two"}, []int{1, 2, 2}},

		/* invalid value treat as default strategy */
		{geko.DuplicatedKeyStrategy(10), []string{"one", "three", "two"}, []int{1, 3, 2}},
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
	ps := geko.NewPairs[int, string]()
	ps.Add(3, "three.2")
	ps.Add(1, "one")
	ps.Add(4, "four")
	ps.Add(2, "two")
	ps.Add(3, "three.1")

	ps.Sort(func(a, b *geko.Pair[int, string]) bool {
		return a.Key < b.Key
	})

	exceptedPairs := []geko.Pair[int, string]{
		{1, "one"},
		{2, "two"},
		{3, "three.2"},
		{3, "three.1"}, // Test sort stability
		{4, "four"},
	}

	if !reflect.DeepEqual(ps.List, exceptedPairs) {
		t.Fatalf("Sort result excepted %#v, got %#v", exceptedPairs, ps.List)
	}
}

func TestPairs_Filter(t *testing.T) {
	ps := geko.NewPairs[int, string]()
	ps.Add(1, "one")
	ps.Add(2, "two")
	ps.Add(3, "three")
	ps.Add(4, "four")

	ps.Filter(func(p *geko.Pair[int, string]) bool {
		return p.Key%2 == 0
	})

	exceptedPairs := []geko.Pair[int, string]{
		{2, "two"},
		{4, "four"},
	}

	if !reflect.DeepEqual(exceptedPairs, ps.List) {
		t.Fatalf("Filter result excepted %#v, got %#v", exceptedPairs, ps.List)
	}
}

func TestPairs_MarshalJSON_InvalidKeyType(t *testing.T) {
	marshalWillReportError[*json.UnsupportedTypeError](t, geko.NewPairs[int, string]())
	marshalWillReportError[*json.UnsupportedTypeError](t, geko.NewPairs[*string, int]())

	type myString string
	marshalWillReportError[*json.UnsupportedTypeError](t, geko.NewPairs[myString, int]())
}

func TestPairs_MarshalJSON_Nil(t *testing.T) {
	var ps *geko.Pairs[string, int]

	data, err := json.Marshal(ps)
	if err != nil {
		t.Fatalf("Marshal nil pairs with error: %s", err.Error())
	}

	if string(data) != `null` {
		t.Fatalf("Marshal result %s not correct", string(data))
	}
}

func TestMap_MarshalJSON_EmptyPairs(t *testing.T) {
	ps := geko.NewPairs[string, any]()

	data, err := json.Marshal(ps)
	if err != nil {
		t.Fatalf("Marshal empty pairs with error: %s", err.Error())
	}

	if string(data) != `{}` {
		t.Fatalf("Marshal result %s not correct", string(data))
	}
}

func TestPairs_MarshalJSON_StringToInt(t *testing.T) {
	ps := geko.NewPairs[string, int]()
	ps.Add("z", 1)
	ps.Add("a", 2)
	ps.Add("n", 3)

	data, err := json.Marshal(ps)
	if err != nil {
		t.Fatalf("Marshal %#v with error: %s", ps, err.Error())
	}

	if string(data) != `{"z":1,"a":2,"n":3}` {
		t.Fatalf("Marshal result %s not correct", string(data))
	}
}

func TestPairs_MarshalJSON_StringToAny(t *testing.T) {
	ps := geko.NewPairs[string, any]()

	ps.Add("string", "hello")
	ps.Add("number", 2)
	ps.Add("float", 2.5)
	ps.Add("json_number", json.Number("10"))
	ps.Add("array", []any{7, "s"})
	ps.Add("bool", true)
	ps.Add("null", nil)

	data, err := json.Marshal(ps)
	if err != nil {
		t.Fatalf("Marshal %#v with error: %s", ps, err.Error())
	}

	if string(data) != `{`+
		`"string":"hello","number":2,"float":2.5,"json_number":10,`+
		`"array":[7,"s"],"bool":true,"null":null`+
		`}` {
		t.Fatalf("Marshal result %s not correct", string(data))
	}
}

func TestPairs_UnmarshalJSON_DirectlyCallWithInvalidData(t *testing.T) {
	ps := geko.NewPairs[string, any]()
	if err := ps.UnmarshalJSON([]byte("")); err == nil {
		t.Fatalf("Should report error with empty input")
	}
	if err := ps.UnmarshalJSON([]byte(`x`)); err == nil {
		t.Fatalf("Should report error with invalid input")
	}
}

func TestPairs_UnmarshalJSON_NilPairs(t *testing.T) {
	var ps geko.ObjectItems
	if err := json.Unmarshal([]byte(`{"a": 1}`), ps); err == nil {
		t.Fatalf("Unmarshal into nil pairs do not error")
	}

	// *Pairs = std map, so this format is better, it supports null like std map
	if err := json.Unmarshal([]byte(`{"a": 1}`), &ps); err != nil {
		t.Fatalf("Unmarshal object into pointer to nil pairs with error: %s", err.Error())
	}

	var m2 = geko.NewPairs[string, any]()
	if err := json.Unmarshal([]byte(`null`), &m2); err != nil {
		t.Fatalf("Unmarshal null into pointer to nil pairs with error: %s", err.Error())
	}
	if m2 != nil {
		t.Fatalf("Unmarshal null into Pairs do not get nil")
	}
}

func TestPairs_UnmarshalJSON_InvalidKeyType(t *testing.T) {
	unmarshalWillReportError[*json.UnmarshalTypeError](t, "{}", geko.NewPairs[int, string]())
	unmarshalWillReportError[*json.UnmarshalTypeError](t, "{}", geko.NewPairs[*string, int]())

	type myString string
	unmarshalWillReportError[*json.UnmarshalTypeError](t, "{}", geko.NewPairs[myString, int]())
}

func TestPairs_UnmarshalJSON_UnmatchedType(t *testing.T) {
	unmarshalWillReportError[*json.UnmarshalTypeError](t, "[]", geko.NewPairs[string, any]())
	unmarshalWillReportError[*json.UnmarshalTypeError](t, "4", geko.NewPairs[string, any]())
	unmarshalWillReportError[*json.UnmarshalTypeError](t, `"string"`, geko.NewPairs[string, any]())
	unmarshalWillReportError[*json.UnmarshalTypeError](t, "true", geko.NewPairs[string, any]())
}

func TestPairs_UnmarshalJSON_UnmatchedValueType(t *testing.T) {
	unmarshalWillReportError[*json.UnmarshalTypeError](t, `{"a":"str"}`, geko.NewPairs[string, int]())
}

func TestPairs_UnmarshalJSON_ConcreteValueType(t *testing.T) {
	ps := geko.NewPairs[string, int]()
	if err := json.Unmarshal([]byte(`{"a": 1}`), &ps); err != nil {
		t.Fatalf("Unmarshal with error: %s", err.Error())
	}

	ps2 := geko.NewPairs[string, s]()
	if err := json.Unmarshal([]byte(`{"a": {"s": "good"}}`), &ps2); err != nil {
		t.Fatalf("Unmarshal with error: %s", err.Error())
	}
	if ps2.GetValueByIndex(0).S != "good" {
		t.Fatalf("Unmarshal into concrete struct type failed: %#v", ps2)
	}

	ps3 := geko.NewPairs[string, *s]()
	if err := json.Unmarshal([]byte(`{"a": {"s": "good"}}`), &ps3); err != nil {
		t.Fatalf("Unmarshal with error: %s", err.Error())
	}
	if ps3.GetValueByIndex(0).S != "good" {
		t.Fatalf("Unmarshal into concrete struct pointer type failed: %#v", ps3)
	}
}

func TestPairs_UnmarshalJSON_InitializedPairs(t *testing.T) {
	ps := geko.NewPairs[string, any]()
	ps.Add("old", "value")
	if err := json.Unmarshal([]byte(`{"a": 1}`), &ps); err != nil {
		t.Fatalf("Unmarshal into initialized map with error: %s", err.Error())
	}

	exceptedKeys := []string{"old", "a"}
	keys := ps.Keys()
	if !reflect.DeepEqual(keys, exceptedKeys) {
		t.Fatalf("Excepted keys %#v, got %#v", exceptedKeys, keys)
	}
}

func TestPairs_UnmarshalJSON_DuplicatedKey(t *testing.T) {
	ps := geko.NewPairs[string, int]()

	data := []byte(`{"a":1,"b":2,"a":3}`)

	if err := json.Unmarshal(data, &ps); err != nil {
		t.Fatalf("Unmarshal error: %s", err.Error())
	}

	excepted := []geko.Pair[string, int]{
		{"a", 1},
		{"b", 2},
		{"a", 3},
	}

	if !reflect.DeepEqual(ps.List, excepted) {
		t.Fatalf(
			"Excepted keys %#v, got %#v",
			excepted, ps.List,
		)
	}

	output, err := json.Marshal(ps)
	if err != nil {
		t.Fatalf("Marshal error: %s", err.Error())
	}

	if !reflect.DeepEqual(output, data) {
		t.Fatalf("Marshal result not same as input: %s", string(output))
	}
}

func TestPairs_UnmarshalJSON_InnerValueUseOurType(t *testing.T) {
	ps := geko.NewPairs[string, any]()
	if err := json.Unmarshal([]byte(`{"arr":[1,2,{"a":1,"b":2,"a":3}]}`), &ps); err != nil {
		t.Fatalf("Unmarshal error: %s", err.Error())
	}

	arr := ps.Get("arr")
	if len(arr) == 0 {
		t.Fatalf("Key arr not exist")
	}

	l, ok := arr[0].(geko.Array)
	if !ok {
		t.Fatalf("Inner array is not List type")
	}

	obj := l.Get(2)

	_, ok = obj.(geko.ObjectItems)
	if !ok {
		t.Fatalf("Inner array -> object is not Pairs type")
	}
}
