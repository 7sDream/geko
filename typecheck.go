package geko

import "reflect"

func isAny[T any]() bool {
	var checker T
	var checkerRef = reflect.ValueOf(&checker).Elem()

	return checkerRef.Kind() == reflect.Interface && checkerRef.NumMethod() == 0
}

func underlyingIsString[T any]() bool {
	var checker T
	var checkerTyp = reflect.TypeOf(checker)

	return checkerTyp != nil && checkerTyp.Kind() == reflect.String
}
