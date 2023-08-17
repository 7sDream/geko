package geko_test

import (
	"encoding/json"
	"errors"
	"reflect"
	"testing"

	"github.com/7sDream/geko"
)

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
		t.Fatalf("Unmarshal *%s type error is not excepted type", typ)
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

func TestMap_UnmarshalJSON_InvalidKeyType(t *testing.T) {
	unmarshalWillReportError[*json.UnmarshalTypeError](t, "{}", geko.NewMap[int, string]())
	unmarshalWillReportError[*json.UnmarshalTypeError](t, "{}", geko.NewMap[*string, int]())

	type myString string
	unmarshalWillReportError[*json.UnmarshalTypeError](t, "{}", geko.NewMap[myString, int]())
}

func TestMap_UnmarshalJSON_InvalidData(t *testing.T) {
	unmarshalWillReportError[*json.UnmarshalTypeError](t, "[]", geko.NewMap[string, any]())
	unmarshalWillReportError[*json.UnmarshalTypeError](t, "4", geko.NewMap[string, any]())
	unmarshalWillReportError[*json.UnmarshalTypeError](t, `"string"`, geko.NewMap[string, any]())
	unmarshalWillReportError[*json.UnmarshalTypeError](t, "true", geko.NewMap[string, any]())
}

func TestMap_UnmarshalJSON_NilMap(t *testing.T) {
	var m *geko.Map[string, any]
	if err := json.Unmarshal([]byte(`{"a": 1}`), m); err == nil {
		t.Fatalf("Unmarshal into nil map do not error")
	}

	// *Map = std map, so this format is better, it supports null like std map
	if err := json.Unmarshal([]byte(`{"a": 1}`), &m); err != nil {
		t.Fatalf("Unmarshal into pointer to nil map with error: %s", err.Error())
	}

	var m2 = geko.NewMap[string, any]()
	if err := json.Unmarshal([]byte(`null`), &m2); err != nil {
		t.Fatalf("Unmarshal into pointer to nil map with error: %s", err.Error())
	}
	if m2 != nil {
		t.Fatalf("Unmarshal null into Map do not get nil")
	}
}

func TestMap_UnmarshalJSON_InitializedMap(t *testing.T) {
	m := geko.NewMap[string, any]()
	m.Add("old", "value")
	if err := json.Unmarshal([]byte(`{"a": 1}`), m); err != nil {
		t.Fatalf("Unmarshal into initialized map with error: %s", err.Error())
	}

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

		if err := json.Unmarshal([]byte(`{"a": 1, "b": 2, "a": 3}`), m); err != nil {
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
		if err := json.Unmarshal([]byte(`{"arr":[1,2,{"a":1,"b":2,"a":3}]}`), m); err != nil {
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

		m2 := obj.(*geko.Map[string, any])
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

func TestList_MarshalJSON_EmptyList(t *testing.T) {
	l := geko.NewList[int]()

	data, err := json.Marshal(l)
	if err != nil {
		t.Fatalf("Marshal empty list with error: %s", err.Error())
	}

	if string(data) != `[]` {
		t.Fatalf("Marshal result %s not correct", string(data))
	}
}
