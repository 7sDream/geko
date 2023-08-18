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

func ExampleMap() {
	m := geko.NewMap[string, int]()

	m.Set("one", 1)
	m.Set("three", 2)
	m.Set("two", 2)
	m.Set("three", 3) // set always do not change order of existed key, so "three" will stay ahead of "two".
	m.Set("four", 0)
	m.Set("five", 5)

	m.SetDuplicatedKeyStrategy(geko.UpdateValueUpdateOrder)
	m.Add("four", 4) // Add will follow DuplicatedKeyStrategy, so now four is last key, and it's value is 4

	for i, length := 0, m.Len(); i < length; i++ {
		pair := m.GetByIndex(i)
		fmt.Printf("%s: %d\n", pair.Key, pair.Value)
	}

	// Output:
	// one: 1
	// three: 3
	// two: 2
	// five: 5
	// four: 4
}

func TestMap_NewWithCapacity(t *testing.T) {
	m := geko.NewMapWithCapacity[string, int](20)

	if reflect.ValueOf(m).Elem().FieldByName("order").Cap() != 20 {
		t.Fatalf("NewMapWithCapacity does not init with capacity")
	}
}

func TestMap_Get(t *testing.T) {
	m := geko.NewMap[string, int]()
	m.Set("one", 1)
	m.Set("two", 2)

	if v, _ := m.Get("one"); v != 1 {
		t.Fatalf("Expect %d, got %d", 1, v)
	}

	if v, _ := m.Get("two"); v != 2 {
		t.Fatalf("Expect %d, got %d", 2, v)
	}

	if _, exist := m.Get("not_exist"); exist != false {
		t.Fatalf("Get a not exist key should return false")
	}

	if m.GetOrZeroValue("not_exist") != 0 {
		t.Fatalf("Get a not exist key should return zero value")
	}

	m2 := geko.NewMap[string, *int]()
	if m2.GetOrZeroValue("not_exist") != nil {
		t.Fatalf("Get a not exist key should return zero value")
	}
}

func TestMap_Has(t *testing.T) {
	m := geko.NewMap[string, int]()
	m.Set("one", 1)
	m.Set("two", 2)

	if !m.Has("one") {
		t.Fatalf("Has said key 'one' does not exist")
	}

	if m.Has("three") {
		t.Fatalf("Has said key 'three' exist")
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
	m := geko.NewMap[string, int]()

	if !willPanic(func() {
		m.GetKeyByIndex(0)
	}) {
		t.Fatalf("GetKeyByIndex with empty map didn't panic")
	}

	m.Set("one", 1)
	m.Set("three", 2)
	m.Set("two", 2)
	m.Set("three", 3)

	if !willPanic(func() {
		m.GetKeyByIndex(-1)
	}) {
		t.Fatalf("GetKeyByIndex negative index didn't panic")
	}

	if !willPanic(func() {
		m.GetKeyByIndex(10)
	}) {
		t.Fatalf("GetKeyByIndex out-of-bound index didn't panic")
	}

	expected := "one"
	if v := m.GetKeyByIndex(0); v != expected {
		t.Fatalf("GetKeyByIndex(0), Expect %#v, got %#v", expected, v)
	}
}

func TestMap_GetByIndex(t *testing.T) {
	m := geko.NewMap[string, int]()

	if !willPanic(func() {
		m.GetByIndex(0)
	}) {
		t.Fatalf("GetByIndex with empty map didn't panic")
	}

	m.Set("one", 1)
	m.Set("three", 2)
	m.Set("two", 2)
	m.Set("three", 3)

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

	expected := geko.Pair[string, int]{Key: "three", Value: 3}
	if v := m.GetByIndex(1); v != expected {
		t.Fatalf("GetByIndex(1), Expect %#v, got %#v", expected, v)
	}
}

func TestMap_GetValueByIndex(t *testing.T) {
	m := geko.NewMap[string, int]()

	if !willPanic(func() {
		m.GetValueByIndex(0)
	}) {
		t.Fatalf("GetValueByIndex with empty map didn't panic")
	}

	m.Set("one", 1)
	m.Set("three", 2)
	m.Set("two", 2)
	m.Set("three", 3)

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

func TestMap_Set(t *testing.T) {
	for _, strategy := range []geko.DuplicatedKeyStrategy{
		geko.UpdateValueKeepOrder,
		geko.UpdateValueUpdateOrder,
		geko.KeepValueUpdateOrder,
		geko.Ignore,
	} {
		m := geko.NewMap[string, int]()
		m.SetDuplicatedKeyStrategy(strategy)
		m.Set("a", 1)
		m.Set("b", 2)
		m.Set("b", 3)
		m.Set("c", 4)
		m.Set("b", 5)

		keys := []string{
			m.GetKeyByIndex(0),
			m.GetKeyByIndex(1),
			m.GetKeyByIndex(2),
		}
		expectedKeys := []string{"a", "b", "c"}
		if !reflect.DeepEqual(keys, expectedKeys) {
			t.Fatalf("After Set, Expect keys %#v, got %#v", expectedKeys, keys)
		}

		values := []int{
			m.GetOrZeroValue("a"),
			m.GetOrZeroValue("b"),
			m.GetOrZeroValue("c"),
		}
		expectedValues := []int{1, 5, 4}
		if !reflect.DeepEqual(values, expectedValues) {
			t.Fatalf("After Set, Expect keys %#v, got %#v", expectedValues, values)
		}
	}
}

func TestMap_Add(t *testing.T) {
	cases := []struct {
		strategy       geko.DuplicatedKeyStrategy
		exceptedKeys   []string
		exceptedValues []int
	}{
		{geko.UpdateValueKeepOrder, []string{"a", "b"}, []int{3, 2}},
		{geko.UpdateValueUpdateOrder, []string{"b", "a"}, []int{2, 3}},
		{geko.KeepValueUpdateOrder, []string{"b", "a"}, []int{2, 1}},
		{geko.Ignore, []string{"a", "b"}, []int{1, 2}},

		/* invalid value treat as default strategy */
		{geko.DuplicatedKeyStrategy(10), []string{"a", "b"}, []int{3, 2}},
	}

	for _, tt := range cases {
		m := geko.NewMap[string, int]()
		m.SetDuplicatedKeyStrategy(tt.strategy)
		m.Add("a", 1)
		m.Add("b", 2)
		m.Add("a", 3)

		if strategy := m.DuplicatedKeyStrategy(); strategy != tt.strategy {
			t.Fatalf(
				"Excepted strategy %#v, got %#v",
				tt.strategy, strategy,
			)
		}

		keys := m.Keys()
		values := m.Values()

		if !reflect.DeepEqual(keys, tt.exceptedKeys) {
			t.Fatalf(
				"For strategy %#v, excepted keys %#v, got %#v",
				tt.strategy, tt.exceptedKeys, keys,
			)
		}

		if !reflect.DeepEqual(values, tt.exceptedValues) {
			t.Fatalf(
				"For strategy %#v, excepted values %#v, got %#v",
				tt.strategy, tt.exceptedValues, values,
			)
		}
	}
}

func TestMap_Append(t *testing.T) {
	m := geko.NewMap[string, int]()
	m.Append([]geko.Pair[string, int]{
		{"s", 2},
		{"z", 7},
		{"z", 4},
		{"w", 9},
		{"z", 1},
	}...)

	keys := []string{
		m.GetKeyByIndex(0),
		m.GetKeyByIndex(1),
		m.GetKeyByIndex(2),
	}
	expectedKeys := []string{"s", "z", "w"}
	if !reflect.DeepEqual(keys, expectedKeys) {
		t.Fatalf("After Append, expect keys %#v, got %#v", expectedKeys, keys)
	}

	values := []int{
		m.GetOrZeroValue("s"),
		m.GetOrZeroValue("z"),
		m.GetOrZeroValue("w"),
	}
	expectedValues := []int{2, 1, 9}
	if !reflect.DeepEqual(values, expectedValues) {
		t.Fatalf("After Append, expect keys %#v, got %#v", expectedValues, values)
	}
}

func TestMap_Delete(t *testing.T) {
	m := geko.NewMap[string, int]()
	m.Set("a", 1)

	m.Delete("b") // should not panic

	m.Delete("a")

	if m.Len() != 0 {
		t.Fatalf("After Delete all item, Map is not empty")
	}

	m = geko.NewMap[string, int]()
	m.Set("a", 1)
	m.Set("b", 2)
	m.Set("c", 3)
	m.Delete("b")

	if m.Len() != 2 {
		t.Fatalf("After Delete a item, Len does not correct")
	}

	if _, exist := m.Get("b"); exist != false {
		t.Fatalf("After Delete item, it still exist")
	}

	m.Set("b", 4)

	if m.Len() != 3 {
		t.Fatalf("After Delete and Set a same key, Len does not correct")
	}

	if v := m.GetValueByIndex(2); v != 4 {
		t.Fatalf("Item does not appear in last after Delete and Set")
	}
}

func TestMap_DeleteByIndex(t *testing.T) {
	m := geko.NewMap[string, int]()

	if !willPanic(func() {
		m.DeleteByIndex(1)
	}) {
		t.Fatalf("DeleteByIndex with empty map didn't panic")
	}

	m.Set("a", 1)
	m.Set("b", 2)
	m.Set("c", 3)
	m.DeleteByIndex(1)

	if m.Len() != 2 {
		t.Fatalf("After DeleteByIndex, Len does not correct")
	}

	if _, exist := m.Get("b"); exist {
		t.Fatalf("After DeleteByIndex, it still exist")
	}

	keys := m.Keys()
	excepted := []string{"a", "c"}
	if !reflect.DeepEqual(keys, excepted) {
		t.Fatalf("After DeleteByIndex, excepted keys %#v, got %#v", excepted, keys)
	}
}

func TestMap_Clear(t *testing.T) {
	m := geko.NewMap[string, int]()
	m.Set("a", 1)
	m.Set("b", 2)
	m.Clear()

	if m.Len() != 0 {
		t.Fatalf("After Clean, map is not empty")
	}

	if len(m.Keys()) != 0 {
		t.Fatalf("After Clean, map Keys() is not empty")
	}

	// After Clear, new Set should not panic
	m.Set("b", 2)
	m.Set("a", 1)
	keys := m.Keys()
	excepted := []string{"b", "a"}
	if !reflect.DeepEqual(keys, excepted) {
		t.Fatalf("After Clean, old values should not effect new order")
	}
}

func TestMap_Len(t *testing.T) {
	for times := 0; times < 20; times++ {
		exceptedLength := rand.Int() % 100

		m := geko.NewMap[string, int]()
		for i := 0; i < exceptedLength; i++ {
			m.Set(strconv.Itoa(i), i)
			// Add some existing key with a certain probability
			// Shouldn't effect length
			if rand.Int()%3 == 0 {
				m.Set(strconv.Itoa(rand.Int()%exceptedLength), rand.Int())
			}
		}

		length := m.Len()
		if length != exceptedLength {
			t.Fatalf("Length excepted %d, got %d", exceptedLength, length)
		}
	}
}

func TestMap_Keys(t *testing.T) {
	m := geko.NewMap[string, int]()
	m.Set("one", 1)
	m.Set("three", 2)
	m.Set("two", 2)
	m.Set("three", 3)

	m.Delete("one")

	excepted := []string{"three", "two"}
	keys := m.Keys()
	if !reflect.DeepEqual(keys, excepted) {
		t.Fatalf("Excepted keys %#v, got %#v", excepted, keys)
	}

	keys[0] = "haha"
	if reflect.DeepEqual(keys, m.Keys()) {
		t.Fatalf("Modify return keys should not effect map")
	}
}

func TestMap_Values(t *testing.T) {
	m := geko.NewMap[string, int]()
	m.Set("one", 1)
	m.Set("three", 2)
	m.Set("two", 2)
	m.Set("three", 3)

	m.Delete("one")

	excepted := []int{3, 2}
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

	m2 := geko.NewMap[string, *s]()
	m2.Set("one", &s{Value: 1})
	m2.Set("two", &s{Value: 2})
	m2.Set("three", &s{Value: 3})

	m2.Values()[2].Value = 100

	if m2.GetOrZeroValue("three").Value != 100 {
		t.Fatalf("Use pointer as value type will allow user modifier inner value")
	}
}

func TestMap_Pairs(t *testing.T) {
	m := geko.NewMap[string, int]()
	m.Set("one", 1)
	m.Set("three", 2)
	m.Set("two", 2)
	m.Set("three", 3)
	m.Delete("one")

	expected := []geko.Pair[string, int]{
		{"three", 3},
		{"two", 2},
	}
	pairs := m.Pairs().List
	if !reflect.DeepEqual(pairs, expected) {
		t.Fatalf("Excepted %#v, got %#v", expected, pairs)
	}
}

func TestMap_Sort(t *testing.T) {
	m := geko.NewMap[int, string]()
	m.Set(3, "three")
	m.Set(1, "one")
	m.Set(4, "four")
	m.Set(2, "two")

	m.Sort(func(a, b *geko.Pair[int, string]) bool {
		return a.Key < b.Key
	})

	exceptedPairs := []geko.Pair[int, string]{
		{1, "one"},
		{2, "two"},
		{3, "three"},
		{4, "four"},
	}

	pairs := m.Pairs().List

	if !reflect.DeepEqual(pairs, exceptedPairs) {
		t.Fatalf("Sort result excepted %#v, got %#v", exceptedPairs, pairs)
	}
}

func TestMap_Filter(t *testing.T) {
	m := geko.NewMap[int, string]()
	m.Set(1, "one")
	m.Set(2, "two")
	m.Set(3, "three")
	m.Set(4, "four")

	m.Filter(func(p *geko.Pair[int, string]) bool {
		return p.Key%2 == 0
	})

	exceptedPairs := []geko.Pair[int, string]{
		{2, "two"},
		{4, "four"},
	}

	pairs := m.Pairs().List

	if !reflect.DeepEqual(pairs, exceptedPairs) {
		t.Fatalf("Filter result excepted %#v, got %#v", exceptedPairs, pairs)
	}
}

func TestMap_MarshalJSON_InvalidKeyType(t *testing.T) {
	marshalWillReportError[*json.UnsupportedTypeError](t, geko.NewMap[int, string]())
	marshalWillReportError[*json.UnsupportedTypeError](t, geko.NewMap[*string, int]())

	type myString string
	marshalWillReportError[*json.UnsupportedTypeError](t, geko.NewMap[myString, int]())
}

func TestMap_MarshalJSON_Nil(t *testing.T) {
	var m *geko.Map[string, int]

	data, err := json.Marshal(m)
	if err != nil {
		t.Fatalf("Marshal nil map with error: %s", err.Error())
	}

	if string(data) != `null` {
		t.Fatalf("Marshal result %s not correct", string(data))
	}
}

func TestMap_MarshalJSON_ValueError(t *testing.T) {
	m := geko.NewMap[string, any]()
	m.Add("invalid", json.Number("invalid"))

	if _, err := json.Marshal(m); err == nil {
		t.Fatalf("Marshal invalid number do not error")
	}
}

func TestMap_MarshalJSON_EmptyMap(t *testing.T) {
	m := geko.NewMap[string, any]()

	data, err := json.Marshal(m)
	if err != nil {
		t.Fatalf("Marshal empty map with error: %s", err.Error())
	}

	if string(data) != `{}` {
		t.Fatalf("Marshal result %s not correct", string(data))
	}
}

func TestMap_MarshalJSON_StringToInt(t *testing.T) {
	m := geko.NewMap[string, int]()
	m.Set("z", 1)
	m.Set("a", 2)
	m.Set("n", 3)

	data, err := json.Marshal(m)
	if err != nil {
		t.Fatalf("Marshal %#v with error: %s", m, err.Error())
	}

	if string(data) != `{"z":1,"a":2,"n":3}` {
		t.Fatalf("Marshal result %s not correct", string(data))
	}
}

func TestMap_MarshalJSON_StringToAny(t *testing.T) {
	mAny := geko.NewMap[string, any]()

	mAny.Set("string", "hello")
	mAny.Set("number", 2)
	mAny.Set("float", 2.5)
	mAny.Set("json_number", json.Number("10"))
	mAny.Set("array", []any{7, "s"})
	mAny.Set("bool", true)
	mAny.Set("null", nil)

	data, err := json.Marshal(mAny)
	if err != nil {
		t.Fatalf("Marshal %#v with error: %s", mAny, err.Error())
	}

	if string(data) != `{`+
		`"string":"hello","number":2,"float":2.5,"json_number":10,`+
		`"array":[7,"s"],"bool":true,"null":null`+
		`}` {
		t.Fatalf("Marshal result %s not correct", string(data))
	}
}

func TestMap_UnmarshalJSON_DirectlyCallWithInvalidData(t *testing.T) {
	m := geko.NewMap[string, any]()
	if err := m.UnmarshalJSON([]byte("")); err == nil {
		t.Fatalf("Should report error with empty input")
	}
	if err := m.UnmarshalJSON([]byte(`x`)); err == nil {
		t.Fatalf("Should report error with invalid input")
	}
}

func TestMap_UnmarshalJSON_NilMap(t *testing.T) {
	var m geko.Object
	if err := json.Unmarshal([]byte(`{"a": 1}`), m); err == nil {
		t.Fatalf("Unmarshal into nil map do not error")
	}

	// *Map = std map, so this format is better, it supports null like std map
	if err := json.Unmarshal([]byte(`{"a": 1}`), &m); err != nil {
		t.Fatalf("Unmarshal object into pointer to nil map with error: %s", err.Error())
	}

	var m2 = geko.NewMap[string, any]()
	if err := json.Unmarshal([]byte(`null`), &m2); err != nil {
		t.Fatalf("Unmarshal null into pointer to nil map with error: %s", err.Error())
	}
	if m2 != nil {
		t.Fatalf("Unmarshal null into Map do not get nil")
	}
}

func TestMap_UnmarshalJSON_InvalidKeyType(t *testing.T) {
	unmarshalWillReportError[*json.UnmarshalTypeError](t, "{}", geko.NewMap[int, string]())
	unmarshalWillReportError[*json.UnmarshalTypeError](t, "{}", geko.NewMap[*string, int]())

	type myString string
	unmarshalWillReportError[*json.UnmarshalTypeError](t, "{}", geko.NewMap[myString, int]())
}

func TestMap_UnmarshalJSON_UnmatchedType(t *testing.T) {
	unmarshalWillReportError[*json.UnmarshalTypeError](t, "[]", geko.NewMap[string, any]())
	unmarshalWillReportError[*json.UnmarshalTypeError](t, "4", geko.NewMap[string, any]())
	unmarshalWillReportError[*json.UnmarshalTypeError](t, `"string"`, geko.NewMap[string, any]())
	unmarshalWillReportError[*json.UnmarshalTypeError](t, "true", geko.NewMap[string, any]())
}

func TestMap_UnmarshalJSON_UnmatchedValueType(t *testing.T) {
	unmarshalWillReportError[*json.UnmarshalTypeError](t, `{"a":"str"}`, geko.NewMap[string, int]())
}

func TestMap_UnmarshalJSON_ConcreteValueType(t *testing.T) {
	m := geko.NewMap[string, int]()
	if err := json.Unmarshal([]byte(`{"a": 1}`), &m); err != nil {
		t.Fatalf("Unmarshal with error: %s", err.Error())
	}

	m2 := geko.NewMap[string, s]()
	if err := json.Unmarshal([]byte(`{"a": {"s": "good"}}`), &m2); err != nil {
		t.Fatalf("Unmarshal with error: %s", err.Error())
	}
	if m2.GetOrZeroValue("a").S != "good" {
		t.Fatalf("Unmarshal into concrete struct type failed: %#v", m2)
	}

	m3 := geko.NewMap[string, *s]()
	if err := json.Unmarshal([]byte(`{"a": {"s": "good"}}`), &m3); err != nil {
		t.Fatalf("Unmarshal with error: %s", err.Error())
	}
	if m3.GetOrZeroValue("a").S != "good" {
		t.Fatalf("Unmarshal into concrete struct pointer type failed: %#v", m3)
	}
}

func TestMap_UnmarshalJSON_InitializedMap(t *testing.T) {
	m := geko.NewMap[string, any]()
	m.Add("old", "value")
	if err := json.Unmarshal([]byte(`{"a": 1}`), &m); err != nil {
		t.Fatalf("Unmarshal into initialized map with error: %s", err.Error())
	}

	// The behavior of the standard library is to retain the old values
	// and we are consistent with it

	exceptedKeys := []string{"old", "a"}
	keys := m.Keys()
	if !reflect.DeepEqual(keys, exceptedKeys) {
		t.Fatalf("Excepted keys %#v, got %#v", exceptedKeys, keys)
	}
}

func TestMap_UnmarshalJSON_DuplicatedKey(t *testing.T) {
	cases := []struct {
		strategy       geko.DuplicatedKeyStrategy
		exceptedKeys   []string
		exceptedValues []int
	}{
		{geko.UpdateValueKeepOrder, []string{"a", "b"}, []int{3, 2}},
		{geko.UpdateValueUpdateOrder, []string{"b", "a"}, []int{2, 3}},
		{geko.KeepValueUpdateOrder, []string{"b", "a"}, []int{2, 1}},
		{geko.Ignore, []string{"a", "b"}, []int{1, 2}},

		/* invalid value treat as default strategy */
		{geko.DuplicatedKeyStrategy(10), []string{"a", "b"}, []int{3, 2}},
	}

	for _, tt := range cases {
		m := geko.NewMap[string, int]()
		m.SetDuplicatedKeyStrategy(tt.strategy)

		if err := json.Unmarshal([]byte(`{"a": 1, "b": 2, "a": 3}`), &m); err != nil {
			t.Fatalf("Strategy %#v, unmarshal error: %s", tt.strategy, err.Error())
		}

		keys := m.Keys()
		values := m.Values()

		if !reflect.DeepEqual(keys, tt.exceptedKeys) {
			t.Fatalf(
				"For strategy %#v, excepted keys %#v, got %#v",
				tt.strategy, tt.exceptedKeys, keys,
			)
		}

		if !reflect.DeepEqual(values, tt.exceptedValues) {
			t.Fatalf(
				"For strategy %#v, excepted values %#v, got %#v",
				tt.strategy, tt.exceptedValues, values,
			)
		}
	}
}

func TestMap_UnmarshalJSON_InnerValueUseOurType(t *testing.T) {
	cases := []struct {
		strategy       geko.DuplicatedKeyStrategy
		exceptedKeys   []string
		exceptedValues []any
	}{
		{geko.UpdateValueKeepOrder, []string{"a", "b"}, []any{3.0, 2.0}},
		{geko.UpdateValueUpdateOrder, []string{"b", "a"}, []any{2.0, 3.0}},
		{geko.KeepValueUpdateOrder, []string{"b", "a"}, []any{2.0, 1.0}},
		{geko.Ignore, []string{"a", "b"}, []any{1.0, 2.0}},
	}
	for _, tt := range cases {
		m := geko.NewMap[string, any]()
		m.SetDuplicatedKeyStrategy(tt.strategy)
		if err := json.Unmarshal([]byte(`{"arr":[1,2,{"a":1,"b":2,"a":3}]}`), &m); err != nil {
			t.Fatalf("Unmarshal error: %s", err.Error())
		}

		arr, exist := m.Get("arr")
		if !exist {
			t.Fatalf("Key arr not exist")
		}

		l, ok := arr.(geko.Array)
		if !ok {
			t.Fatalf("Inner array is not List type")
		}

		obj := l.Get(2)

		m2, ok := obj.(geko.Object)
		if !ok {
			t.Fatalf("Inner object is not Map type")
		}

		keys := m2.Keys()
		if !reflect.DeepEqual(keys, tt.exceptedKeys) {
			t.Fatalf(
				"Inner object do not follow strategy %#v, excepted keys %#v, got %#v",
				tt.strategy, tt.exceptedKeys, keys,
			)
		}

		values := m2.Values()
		if !reflect.DeepEqual(values, tt.exceptedValues) {
			t.Fatalf(
				"Inner object do not follow strategy %#v, excepted values %#v, got %#v",
				tt.strategy, tt.exceptedValues, values,
			)
		}
	}
}
