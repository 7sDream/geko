package geko

import (
	"testing"
)

type emptyInterface interface {
}

type myString string

func TestIsAny(t *testing.T) {
	if isAny[string]() {
		t.Fatalf("isAny failed in type string")
	}

	if isAny[int]() {
		t.Fatalf("isAny failed in type string")
	}

	if !isAny[interface{}]() {
		t.Fatalf("isAny failed in type interface{}")
	}

	if !isAny[emptyInterface]() {
		t.Fatalf("isAny failed in type EmptyInterface")
	}

	if !isAny[any]() {
		t.Fatalf("isAny failed in type EmptyInterface")
	}
}

func TestIsString(t *testing.T) {
	if !underlyingIsString[string]() {
		t.Fatalf("isString failed in type string")
	}

	if !underlyingIsString[myString]() {
		t.Fatalf("isString failed in type myString")
	}

	if underlyingIsString[int]() {
		t.Fatalf("isString failed in type int")
	}

	if underlyingIsString[interface{}]() {
		t.Fatalf("isString failed in type interface{}")
	}

	if underlyingIsString[emptyInterface]() {
		t.Fatalf("isString failed in type EmptyInterface")
	}

	if underlyingIsString[any]() {
		t.Fatalf("isString failed in type EmptyInterface")
	}
}
