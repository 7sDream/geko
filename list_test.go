package geko_test

import (
	"encoding/json"
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

func TestList_MarshalJSON_Nil(t *testing.T) {
	var l *geko.List[int]

	output, err := json.Marshal(l)
	if err != nil {
		t.Fatalf("Marshal nil list with error: %s", err.Error())
	}

	if string(output) != `null` {
		t.Fatalf("Marshal result %s not correct", string(output))
	}
}

func TestList_MarshalJSON_InternalNilList(t *testing.T) {
	l := geko.NewList[int]()

	output, err := json.Marshal(l)
	if err != nil {
		t.Fatalf("Marshal empty list with error: %s", err.Error())
	}

	if string(output) != `[]` {
		t.Fatalf("Marshal result %s not correct", string(output))
	}
}

func TestList_MarshalJSON_ConcreteType(t *testing.T) {
	l := geko.NewListFrom[int]([]int{1, 2, 3})

	output, err := json.Marshal(l)
	if err != nil {
		t.Fatalf("Marshal empty list with error: %s", err.Error())
	}

	if string(output) != `[1,2,3]` {
		t.Fatalf("Marshal result %s not correct", string(output))
	}
}

func TestList_MarshalJSON_AnyType(t *testing.T) {
	l := geko.NewListFrom[any]([]any{
		1, 2.5, true, nil,
		map[string]int{"a": 1},
		geko.NewMap[string, any](),
	})

	output, err := json.Marshal(l)
	if err != nil {
		t.Fatalf("Marshal empty list with error: %s", err.Error())
	}

	if string(output) != `[1,2.5,true,null,{"a":1},{}]` {
		t.Fatalf("Marshal result %s not correct", string(output))
	}
}

func TestList_UnmarshalJSON_DirectlyCallWithInvalidData(t *testing.T) {
	l := geko.NewList[any]()
	if err := l.UnmarshalJSON([]byte("")); err == nil {
		t.Fatalf("Should report error with empty input")
	}
	if err := l.UnmarshalJSON([]byte(`x`)); err == nil {
		t.Fatalf("Should report error with invalid input")
	}
}

func TestList_UnmarshalJSON_NilList(t *testing.T) {
	var l geko.Array
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

	l2 := geko.NewList[s]()
	if err := json.Unmarshal([]byte(`[{"s": "good"}]`), &l2); err != nil {
		t.Fatalf("Unmarshal with error: %s", err.Error())
	}
	if l2.Get(0).S != "good" {
		t.Fatalf("Unmarshal into concrete struct type failed: %#v", l2)
	}

	l3 := geko.NewList[*s]()
	if err := json.Unmarshal([]byte(`[{"s": "good"}]`), &l3); err != nil {
		t.Fatalf("Unmarshal with error: %s", err.Error())
	}
	if l3.Get(0).S != "good" {
		t.Fatalf("Unmarshal into concrete struct pointer type failed: %#v", l3)
	}
}

func TestList_UnmarshalJSON_InitializedList(t *testing.T) {
	l := geko.NewListFrom[int]([]int{7})

	if err := json.Unmarshal([]byte(`[1]`), &l); err != nil {
		t.Fatalf("Unmarshal into initialized map with error: %s", err.Error())
	}

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
	ll, ok := arr.(geko.Array)
	if !ok {
		t.Fatalf("Inner array is not Array type")
	}

	arrObj := ll.Get(1)
	llm, ok := arrObj.(geko.ObjectItems)
	if !ok {
		t.Fatalf("Inner array -> object is not ObjectItems type")
	}

	llmItem := llm.GetByIndex(0)
	if llmItem.Key != "llm" && llmItem.Value != true {
		t.Fatalf("Inner array -> object item not correct: %#v", llm)
	}

	obj := l.Get(2)
	lm, ok := obj.(geko.ObjectItems)
	if !ok {
		t.Fatalf("Inner object is not ObjectItems type")
	}

	objArr := lm.GetValueByIndex(1)
	lml, ok := objArr.(geko.Array)
	if !ok {
		t.Fatalf("Inner object -> array is not Array type")
	}

	if lml.Get(0) != "lml" {
		t.Fatalf("Inner object -> array item not correct: %#v", lml)
	}
}
