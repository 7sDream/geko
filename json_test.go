package geko_test

import (
	"encoding/json"
	"errors"
	"reflect"
	"testing"

	"github.com/7sDream/geko"
)

type s struct {
	S string `json:"s"`
}

func marshalWillReportError[T error](t *testing.T, v any) {
	_, err := json.Marshal(v)
	if err == nil {
		t.Fatalf("Marshal %s type without error", reflect.TypeOf(v).Name())
	}
	var typeErr T
	if !errors.As(err, &typeErr) {
		t.Fatalf(
			"Marshal %s type error is not %s",
			reflect.TypeOf(v).Elem().Name(),
			reflect.TypeOf(typeErr).Elem().Name(),
		)
	}
}

func unmarshalWillReportError[T error](t *testing.T, data string, v any) {
	typ := reflect.TypeOf(v).Elem().Name()
	err := json.Unmarshal([]byte(data), &v)
	if err == nil {
		t.Fatalf("Unmarshal *%s type without error", typ)
	}
	var typeErr T
	if !errors.As(err, &typeErr) {
		t.Logf("Unmarshal error: %s", err.Error())
		t.Fatalf(
			"Unmarshal *%s type error is not excepted type, it's %s", typ,
			reflect.TypeOf(err).Elem().Name(),
		)
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
	var m *geko.Map[string, any]
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

		l, ok := arr.(*geko.List[any])
		if !ok {
			t.Fatalf("Inner array is not List type")
		}

		obj := l.Get(2)

		m2, ok := obj.(*geko.Map[string, any])
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

func TestList_MarshalJSON_Nil(t *testing.T) {
	var l *geko.List[int]

	data, err := json.Marshal(l)
	if err != nil {
		t.Fatalf("Marshal nil list with error: %s", err.Error())
	}

	if string(data) != `null` {
		t.Fatalf("Marshal result %s not correct", string(data))
	}
}

func TestList_MarshalJSON_InternalNilList(t *testing.T) {
	l := geko.NewList[int]()

	data, err := json.Marshal(l)
	if err != nil {
		t.Fatalf("Marshal empty list with error: %s", err.Error())
	}

	if string(data) != `[]` {
		t.Fatalf("Marshal result %s not correct", string(data))
	}
}

func TestList_MarshalJSON_ConcreteType(t *testing.T) {
	l := geko.NewListFrom[int]([]int{1, 2, 3})

	data, err := json.Marshal(l)
	if err != nil {
		t.Fatalf("Marshal empty list with error: %s", err.Error())
	}

	if string(data) != `[1,2,3]` {
		t.Fatalf("Marshal result %s not correct", string(data))
	}
}

func TestList_MarshalJSON_AnyType(t *testing.T) {
	l := geko.NewListFrom[any]([]any{
		1, 2.5, true, nil,
		map[string]int{"a": 1},
		geko.NewMap[string, any](),
	})

	data, err := json.Marshal(l)
	if err != nil {
		t.Fatalf("Marshal empty list with error: %s", err.Error())
	}

	if string(data) != `[1,2.5,true,null,{"a":1},{}]` {
		t.Fatalf("Marshal result %s not correct", string(data))
	}
}

func TestList_UnmarshalJSON_DirectlyCallWithInvalidData(t *testing.T) {
	m := geko.NewList[any]()
	if err := m.UnmarshalJSON([]byte("")); err == nil {
		t.Fatalf("Should report error with empty input")
	}
	if err := m.UnmarshalJSON([]byte(`x`)); err == nil {
		t.Fatalf("Should report error with invalid input")
	}
}

func TestList_UnmarshalJSON_NilList(t *testing.T) {
	var l *geko.List[any]
	if err := json.Unmarshal([]byte(`[1]`), l); err == nil {
		t.Fatalf("Unmarshal into nil list should report error")
	}

	// *List = std slice, so this format is better, it supports null
	if err := json.Unmarshal([]byte(`[1]`), &l); err != nil {
		t.Fatalf("Unmarshal into pointer to nil list with error: %s", err.Error())
	}

	var m2 = geko.NewList[any]()
	if err := json.Unmarshal([]byte(`null`), &m2); err != nil {
		t.Fatalf("Unmarshal into pointer to nil list with error: %s", err.Error())
	}
	if m2 != nil {
		t.Fatalf("Unmarshal null into Map do not get nil")
	}
}

func TestList_UnmarshalJSON_UnmatchedType(t *testing.T) {
	unmarshalWillReportError[*json.UnmarshalTypeError](t, "{}", geko.NewList[any]())
	unmarshalWillReportError[*json.UnmarshalTypeError](t, "4", geko.NewList[any]())
	unmarshalWillReportError[*json.UnmarshalTypeError](t, `"string"`, geko.NewList[any]())
	unmarshalWillReportError[*json.UnmarshalTypeError](t, "true", geko.NewList[any]())
}

func TestList_UnmarshalJSON_UnmatchedValueType(t *testing.T) {
	unmarshalWillReportError[*json.UnmarshalTypeError](t, `["str"]`, geko.NewList[int]())
	unmarshalWillReportError[*json.UnmarshalTypeError](t, `["str"]`, geko.NewList[s]())
}

func TestList_UnmarshalJSON_ConcreteType(t *testing.T) {
	l := geko.NewList[int]()
	if err := json.Unmarshal([]byte(`[1]`), &l); err != nil {
		t.Fatalf("Unmarshal with error: %s", err.Error())
	}

	m2 := geko.NewList[s]()
	if err := json.Unmarshal([]byte(`[{"s": "good"}]`), &m2); err != nil {
		t.Fatalf("Unmarshal with error: %s", err.Error())
	}
	if m2.Get(0).S != "good" {
		t.Fatalf("Unmarshal into concrete struct type failed: %#v", m2)
	}

	m3 := geko.NewList[*s]()
	if err := json.Unmarshal([]byte(`[{"s": "good"}]`), &m3); err != nil {
		t.Fatalf("Unmarshal with error: %s", err.Error())
	}
	if m3.Get(0).S != "good" {
		t.Fatalf("Unmarshal into concrete struct pointer type failed: %#v", m3)
	}
}

func TestList_UnmarshalJSON_InitializedList(t *testing.T) {
	l := geko.NewListFrom[int]([]int{7})

	if err := json.Unmarshal([]byte(`[1]`), &l); err != nil {
		t.Fatalf("Unmarshal into initialized map with error: %s", err.Error())
	}

	// The behavior of the standard library is to clear the list
	// and we are consistent with it

	excepted := []int{1}
	if !reflect.DeepEqual(l.List, excepted) {
		t.Fatalf("Excepted %#v, got %#v", excepted, l.List)
	}

	l2 := geko.NewListFrom[any]([]any{"old"})
	if err := json.Unmarshal([]byte(`[1]`), &l2); err != nil {
		t.Fatalf("Unmarshal into initialized map with error: %s", err.Error())
	}

	excepted2 := []any{1.0}
	if !reflect.DeepEqual(l2.List, excepted2) {
		t.Fatalf("Excepted %#v, got %#v", excepted2, l2.List)
	}
}

func TestList_UnmarshalJSON_InnerValueUseOurType(t *testing.T) {
	l := geko.NewList[any]()
	if err := json.Unmarshal([]byte(`[1,["2",{"llm":true}],{"a":1,"arr":["lml"]}]`), &l); err != nil {
		t.Fatalf("Unmarshal error: %s", err.Error())
	}

	arr := l.Get(1)
	ll, ok := arr.(*geko.List[any])
	if !ok {
		t.Fatalf("Inner array is not List type")
	}

	arrObj := ll.Get(1)
	llm, ok := arrObj.(*geko.Map[string, any])
	if !ok {
		t.Fatalf("Inner array -> object is not Map type")
	}

	llmItem := llm.GetByIndex(0)
	if llmItem.Key != "llm" && llmItem.Value != true {
		t.Fatalf("Inner array -> object item not correct: %#v", llm)
	}

	obj := l.Get(2)
	lm, ok := obj.(*geko.Map[string, any])
	if !ok {
		t.Fatalf("Inner object is not Map type")
	}

	objArr := lm.GetValueByIndex(1)
	lml, ok := objArr.(*geko.List[any])
	if !ok {
		t.Fatalf("Inner object -> array is not List type")
	}

	if lml.Get(0) != "lml" {
		t.Fatalf("Inner object -> array item not correct: %#v", lml)
	}
}

func TestPairs_MarshalJSON_InvalidKeyType(t *testing.T) {
	marshalWillReportError[*json.UnsupportedTypeError](t, geko.NewPairs[int, string]())
	marshalWillReportError[*json.UnsupportedTypeError](t, geko.NewPairs[*string, int]())

	type myString string
	marshalWillReportError[*json.UnsupportedTypeError](t, geko.NewPairs[myString, int]())
}

func TestPairs_MarshalJSON_Nil(t *testing.T) {
	var m *geko.Pairs[string, int]

	data, err := json.Marshal(m)
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
	psAny := geko.NewPairs[string, any]()

	psAny.Add("string", "hello")
	psAny.Add("number", 2)
	psAny.Add("float", 2.5)
	psAny.Add("json_number", json.Number("10"))
	psAny.Add("array", []any{7, "s"})
	psAny.Add("bool", true)
	psAny.Add("null", nil)

	data, err := json.Marshal(psAny)
	if err != nil {
		t.Fatalf("Marshal %#v with error: %s", psAny, err.Error())
	}

	if string(data) != `{`+
		`"string":"hello","number":2,"float":2.5,"json_number":10,`+
		`"array":[7,"s"],"bool":true,"null":null`+
		`}` {
		t.Fatalf("Marshal result %s not correct", string(data))
	}
}

func TestPairs_UnmarshalJSON_DirectlyCallWithInvalidData(t *testing.T) {
	m := geko.NewPairs[string, any]()
	if err := m.UnmarshalJSON([]byte("")); err == nil {
		t.Fatalf("Should report error with empty input")
	}
	if err := m.UnmarshalJSON([]byte(`x`)); err == nil {
		t.Fatalf("Should report error with invalid input")
	}
}

func TestPairs_UnmarshalJSON_NilPairs(t *testing.T) {
	var m *geko.Pairs[string, any]
	if err := json.Unmarshal([]byte(`{"a": 1}`), m); err == nil {
		t.Fatalf("Unmarshal into nil pairs do not error")
	}

	// *Pairs = std map, so this format is better, it supports null like std map
	if err := json.Unmarshal([]byte(`{"a": 1}`), &m); err != nil {
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
	m := geko.NewPairs[string, any]()
	m.Add("old", "value")
	if err := json.Unmarshal([]byte(`{"a": 1}`), &m); err != nil {
		t.Fatalf("Unmarshal into initialized map with error: %s", err.Error())
	}

	exceptedKeys := []string{"old", "a"}
	keys := m.Keys()
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
	m := geko.NewPairs[string, any]()
	if err := json.Unmarshal([]byte(`{"arr":[1,2,{"a":1,"b":2,"a":3}]}`), &m); err != nil {
		t.Fatalf("Unmarshal error: %s", err.Error())
	}

	arr := m.Get("arr")
	if len(arr) == 0 {
		t.Fatalf("Key arr not exist")
	}

	l, ok := arr[0].(*geko.List[any])
	if !ok {
		t.Fatalf("Inner array is not List type")
	}

	obj := l.Get(2)

	_, ok = obj.(*geko.Pairs[string, any])
	if !ok {
		t.Fatalf("Inner array -> object is not Pairs type")
	}
}

func TestJSONUnmarshal_InvalidData(t *testing.T) {
	invalid := func(data string) {
		_, err := geko.JSONUnmarshal([]byte(data))
		if err == nil {
			t.Fatalf("Do not error with invalid data %s", data)
		}
	}

	invalid(`ss`)
	invalid(`1.34dee`)
	invalid(`[1,2`)
	invalid(`[1,2`)
	invalid(`[1,2,{`)
	invalid(`{s: 1}`)
	invalid(`{"s: 1}`)
	invalid(`{"s": 1,}`)
	invalid(`{"s": 1, "b": {s}}`)
	invalid(`{"s": 1,}ee`)
}

func TestJSONUnmarshal_UseNumber(t *testing.T) {
	data := []byte(`1.245555532383454092038423904345445098423984`)

	result, err := geko.JSONUnmarshal(data, geko.UseNumber(true))
	if err != nil {
		t.Fatalf("Unmarshal error: %s", err.Error())
	}

	if _, ok := result.(json.Number); !ok {
		t.Fatalf("result type is not json.Number: %#v", result)
	}

	data2 := []byte(`{` +
		`"float": 1.245555532383454092038423904345445098423984,` +
		`"int": 245555532383454092038423904345445098423984,` +
		`"arr": [1.2]` +
		`}`)

	result2, err := geko.JSONUnmarshal(data2, geko.UseNumber(true))
	if err != nil {
		t.Fatalf("Unmarshal error: %s", err.Error())
	}

	m, ok := result2.(*geko.Map[string, any])
	if !ok {
		t.Fatalf("result type is not Map: %#v", result)
	}

	f := m.GetOrZeroValue("float")
	if _, ok := f.(json.Number); !ok {
		t.Fatalf("float number in object is not json.Number: %#v", f)
	}

	i := m.GetOrZeroValue("int")
	if _, ok := i.(json.Number); !ok {
		t.Fatalf("int number in object is not json.Number: %#v", i)
	}

	f2 := m.GetOrZeroValue("arr").(*geko.List[any]).Get(0)
	if _, ok := f2.(json.Number); !ok {
		t.Fatalf("float number in object -> list is not json.Number: %#v", f2)
	}
}

func TestJSONUnmarshal_UsePairList(t *testing.T) {
	data := []byte(`{"a":1,"a":2,"obj":{"b":1,"b":2},"arr":[{"c":1,"c":2}]}`)

	result, err := geko.JSONUnmarshal(data, geko.UsePairs(true))
	if err != nil {
		t.Fatalf("Unmarshal error: %s", err.Error())
	}

	ps, ok := result.(*geko.Pairs[string, any])
	if !ok {
		t.Fatalf("outmost object type is not Pairs: %#v", ps)
	}

	obj := ps.Get("obj")[0]
	_, ok = obj.(*geko.Pairs[string, any])
	if !ok {
		t.Fatalf("nets object type is not Pairs: %#v", obj)
	}

	arr := ps.Get("arr")[0]
	arrL, ok := arr.(*geko.List[any])
	if !ok {
		t.Fatalf("inner array type is not List: %#v", arr)
	}

	arrObj := arrL.Get(0)
	_, ok = arrObj.(*geko.Pairs[string, any])
	if !ok {
		t.Fatalf("nets array -> object type is not Pairs: %#v", obj)
	}

	output, err := json.Marshal(ps)
	if err != nil {
		t.Fatalf("Marshal error: %s", err.Error())
	}

	if !reflect.DeepEqual(output, data) {
		t.Fatalf("Marshal result not same as input: %s", string(output))
	}
}

func TestJSONUnmarshal_OnDuplicatedKey(t *testing.T) {
	cases := []struct {
		strategy       geko.DuplicatedKeyStrategy
		exceptedKeys   []string
		exceptedValues []any
	}{
		{geko.UpdateValueKeepOrder, []string{"a", "b"}, []any{3.0, 2.0}},
		{geko.UpdateValueUpdateOrder, []string{"b", "a"}, []any{2.0, 3.0}},
		{geko.KeepValueUpdateOrder, []string{"b", "a"}, []any{2.0, 1.0}},
		{geko.Ignore, []string{"a", "b"}, []any{1.0, 2.0}},

		/* invalid value treat as default strategy */
		{geko.DuplicatedKeyStrategy(10), []string{"a", "b"}, []any{3.0, 2.0}},
	}

	for _, tt := range cases {
		arr, err := geko.JSONUnmarshal(
			[]byte(`[1,false,{"a": 1, "b": 2, "a": 3}]`),
			geko.OnDuplicatedKey(tt.strategy),
		)
		if err != nil {
			t.Fatalf("Strategy %#v, unmarshal error: %s", tt.strategy, err.Error())
		}

		l, _ := arr.(*geko.List[any])

		m := l.Get(2).(*geko.Map[string, any])

		strategy := m.DuplicatedKeyStrategy()
		if strategy != tt.strategy {
			t.Fatalf(
				"For strategy %#v, got strategy %#v",
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
