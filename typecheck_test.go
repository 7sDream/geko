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

type emptyInterface interface {
}

type publicInterface interface {
	Good()
}

type privateInterface interface {
	good()
}

func TestIsAny(t *testing.T) {
	if !isAny[interface{}]() {
		t.Fatalf("isAny failed in type interface{}")
	}

	if !isAny[any]() {
		t.Fatalf("isAny failed in type any")
	}

	if !isAny[emptyInterface]() {
		t.Fatalf("isAny failed in type emptyInterface")
	}

	if isAny[string]() {
		t.Fatalf("isAny failed in type string")
	}

	if isAny[int]() {
		t.Fatalf("isAny failed in type int")
	}

	if isAny[publicInterface]() {
		t.Fatalf("isAny failed in type publicInterface")
	}

	if isAny[privateInterface]() {
		t.Fatalf("isAny failed in type publicInterface")
	}
}

func TestIsString(t *testing.T) {
	if !IsString[string]() {
		t.Fatalf("isString failed in type string")
	}

	if IsString[*string]() {
		t.Fatalf("isString failed in type *string")
	}

	if IsString[int]() {
		t.Fatalf("isString failed in type int")
	}

	if IsString[float64]() {
		t.Fatalf("isString failed in type float64")
	}

	if IsString[bool]() {
		t.Fatalf("isString failed in type bool")
	}

	if IsString[stringWrapper]() {
		t.Fatalf("isString failed in type stringWrapper")
	}

	if IsString[stringStruct]() {
		t.Fatalf("isString failed in type stringStruct")
	}

	if IsString[stringFieldStruct]() {
		t.Fatalf("isString failed in type stringFieldStruct")
	}

	if IsString[realStruct]() {
		t.Fatalf("isString failed in type realStruct")
	}
}
