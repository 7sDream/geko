package geko_test

import (
	"encoding/json"
	"errors"
	"reflect"
	"testing"
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
