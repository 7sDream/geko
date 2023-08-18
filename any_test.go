package geko_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/7sDream/geko"
)

func TestAny_MarshalJSON(t *testing.T) {
	a := geko.Any{}
	data, err := json.Marshal(a)
	if err != nil {
		t.Fatalf("Marshal error: %s", err.Error())
	}
	if string(data) != "null" {
		t.Fatalf("Marshal result not correct: %s", string(data))
	}
}

func TestAny_UnmarshalJSON(t *testing.T) {
	a := geko.Any{}
	err := json.Unmarshal([]byte("null"), &a)
	if err != nil {
		t.Fatalf("Unmarshal error: %s", err.Error())
	}
	if a.Value != nil {
		t.Fatalf("Unmarshal result not correct: %#v", a.Value)
	}
}

func TestAny_UnmarshalJSON_InvalidData(t *testing.T) {
	invalid := func(data string) {
		a := geko.Any{}
		err := a.UnmarshalJSON([]byte(data))
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

	ps, ok := result2.(geko.ObjectItems)
	if !ok {
		t.Fatalf("result type is not Map: %#v", result)
	}

	f := ps.GetFirstOrZeroValue("float")
	if _, ok := f.(json.Number); !ok {
		t.Fatalf("float number in object is not json.Number: %#v", f)
	}

	i := ps.GetFirstOrZeroValue("int")
	if _, ok := i.(json.Number); !ok {
		t.Fatalf("int number in object is not json.Number: %#v", i)
	}

	f2 := ps.GetFirstOrZeroValue("arr").(geko.Array).Get(0)
	if _, ok := f2.(json.Number); !ok {
		t.Fatalf("float number in object -> list is not json.Number: %#v", f2)
	}
}

func TestJSONUnmarshal_UseObjectItem(t *testing.T) {
	data := []byte(`{"a":1,"a":2,"obj":{"b":1,"b":2},"arr":[{"c":1,"c":2}]}`)

	result, err := geko.JSONUnmarshal(data, geko.UseObjectItem())
	if err != nil {
		t.Fatalf("Unmarshal error: %s", err.Error())
	}

	ps, ok := result.(geko.ObjectItems)
	if !ok {
		t.Fatalf("outmost object type is not Pairs: %#v", ps)
	}

	obj := ps.Get("obj")[0]
	_, ok = obj.(geko.ObjectItems)
	if !ok {
		t.Fatalf("nets object type is not Pairs: %#v", obj)
	}

	arr := ps.Get("arr")[0]
	arrL, ok := arr.(geko.Array)
	if !ok {
		t.Fatalf("inner array type is not List: %#v", arr)
	}

	arrObj := arrL.Get(0)
	_, ok = arrObj.(geko.ObjectItems)
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
			geko.UseObject(),
			geko.ObjectOnDuplicatedKey(tt.strategy),
		)
		if err != nil {
			t.Fatalf("Strategy %#v, unmarshal error: %s", tt.strategy, err.Error())
		}

		l, _ := arr.(geko.Array)

		m := l.Get(2).(geko.Object)

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
