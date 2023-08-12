package geko_test

import (
	"encoding/json"
	"testing"

	"github.com/7sDream/geko"
)

func TestJSONUnmarshal(t *testing.T) {
	data := `{"two":2,"one":1,"three":null}`
	value, err := geko.JSONUnmarshal([]byte(data), geko.UseNumber(true))
	if err != nil {
		t.Fatalf("Error: %s", err.Error())
	}

	output, _ := json.Marshal(value)
	t.Logf("marshal result: %s", string(output))
	if string(output) != data {
		t.Fatalf("want %s, got %s", data, string(output))
	}
}

func TestMapUnmarshal(t *testing.T) {
	type myString string
	data := `{"two":2,"one":1,"three":3}`
	kom := geko.NewMap[myString, int]()
	if err := json.Unmarshal([]byte(data), kom); err != nil {
		t.Fatalf("Error: %s", err.Error())
	}

	output, _ := json.Marshal(kom)
	t.Logf("marshal result: %s", string(output))
	if string(output) != data {
		t.Fatalf("want %s, got %s", data, string(output))
	}
}

func TestListUnmarshal(t *testing.T) {
	data := "[3]"
	kol := geko.NewList[any]()
	err := json.Unmarshal([]byte(data), &kol)
	if err != nil {
		t.Fatalf("Error: %s", err.Error())
	}

	output, _ := json.Marshal(kol)
	t.Logf("marshal result: %s", string(output))
	if string(output) != data {
		t.Fatalf("want %s, got %s", data, string(output))
	}
}

func TestMapUnmarshalNestMap(t *testing.T) {
	data := `{"two":2,"one":1,"three":{"five":5,"four":4}}`
	kom := geko.NewMap[string, any]()
	if err := json.Unmarshal([]byte(data), kom); err != nil {
		t.Fatalf("Error: %s", err.Error())
	}

	output, _ := json.Marshal(kom)
	t.Logf("marshal result: %s", string(output))
	if string(output) != data {
		t.Fatalf("want %s, got %s", data, string(output))
	}
}

func TestMapUnmarshalNestArray(t *testing.T) {
	data := `{"two":2,"one":1,"three":["four",4,{"six":6,"five":5}]}`
	kom := geko.NewMap[string, any]()
	if err := json.Unmarshal([]byte(data), kom); err != nil {
		t.Fatalf("Error: %s", err.Error())
	}

	output, _ := json.Marshal(kom)
	t.Logf("marshal result: %s", string(output))
	if string(output) != data {
		t.Fatalf("want %s, got %s", data, string(output))
	}
}