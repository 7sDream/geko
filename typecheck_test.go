package geko

import (
	"testing"
)

type stringWrapper string

type stringStruct struct {
	string
}

type stringFieldStruct struct {
	s string
}

type realStruct struct {
	s string
	i int
}

//nolint:revive
type emptyInterface interface {
}

type publicInterface interface {
	Good()
}

type privateInterface interface {
	good()
}

func TestIsAny(t *testing.T) {
	if !isEmptyInterface[interface{}]() { //nolint:revive
		t.Fatalf("isAny failed in type interface{}")
	}

	if !isEmptyInterface[any]() {
		t.Fatalf("isAny failed in type any")
	}

	if !isEmptyInterface[emptyInterface]() {
		t.Fatalf("isAny failed in type emptyInterface")
	}

	if isEmptyInterface[string]() {
		t.Fatalf("isAny failed in type string")
	}

	if isEmptyInterface[int]() {
		t.Fatalf("isAny failed in type int")
	}

	if isEmptyInterface[publicInterface]() {
		t.Fatalf("isAny failed in type publicInterface")
	}

	if isEmptyInterface[privateInterface]() {
		t.Fatalf("isAny failed in type publicInterface")
	}
}

func TestIsString(t *testing.T) {
	if !isString[string]() {
		t.Fatalf("isString failed in type string")
	}

	if isString[*string]() {
		t.Fatalf("isString failed in type *string")
	}

	if isString[int]() {
		t.Fatalf("isString failed in type int")
	}

	if isString[float64]() {
		t.Fatalf("isString failed in type float64")
	}

	if isString[bool]() {
		t.Fatalf("isString failed in type bool")
	}

	if isString[stringWrapper]() {
		t.Fatalf("isString failed in type stringWrapper")
	}

	if isString[stringStruct]() {
		t.Fatalf("isString failed in type stringStruct")
	}

	if isString[stringFieldStruct]() {
		t.Fatalf("isString failed in type stringFieldStruct")
	}

	if isString[realStruct]() {
		t.Fatalf("isString failed in type realStruct")
	}
}
