package geko

import (
	"bytes"
	"encoding/json"
	"io"
	"reflect"
	"unsafe"
)

// ===== Decoder =====

// JSONUnmarshal likes json.Unmarshal, but uses our [Map]/[Pairs] and [List]
// when meet JSON object and array.
//
// So the type of returned value can be:
// bool, float64/[json.Number], string, nil,
// *[Map][string, any]/*[Pairs][string, any], *[List][any].
//
// The any type in the above container can only be the above type, recursive.
func JSONUnmarshal(data []byte, option ...DecodeOption) (any, error) {
	return newDecoder(bytes.NewReader(data), option...).decode()
}

type decodeOptions struct {
	useNumber             bool
	usePairs              bool
	duplicatedKeyStrategy DuplicatedKeyStrategy
}

// DecodeOption is option type for [JSONUnmarshal].
type DecodeOption func(opts *decodeOptions)

// UseNumber option will make [JSONUnmarshal] uses [json.Number] to store
// JSON number, instead of float64.
func UseNumber(v bool) DecodeOption {
	return func(opts *decodeOptions) {
		opts.useNumber = v
	}
}

// UsePairs option will make [JSONUnmarshal] uses *[Pairs][string, any] to
// store JSON object, instead of *[Map][string, any].
func UsePairs(v bool) DecodeOption {
	return func(opts *decodeOptions) {
		opts.usePairs = v
	}
}

// OnDuplicatedKey set the strategy when there are duplicated key in JSON
// object.
//
// See document of [DuplicatedKeyStrategy] and its enum value for details
func OnDuplicatedKey(strategy DuplicatedKeyStrategy) DecodeOption {
	return func(opts *decodeOptions) {
		opts.duplicatedKeyStrategy = strategy
	}
}

type decoder struct {
	decoder *json.Decoder
	opts    decodeOptions
}

func newDecoder(r io.Reader, option ...DecodeOption) *decoder {
	decoder := &decoder{
		decoder: json.NewDecoder(r),
	}

	for _, opt := range option {
		opt(&decoder.opts)
	}

	return decoder
}

func (d *decoder) decode() (any, error) {
	if d.opts.useNumber {
		d.decoder.UseNumber()
	}

	item, err := d.next()
	if err != nil {
		return nil, err
	}

	if _, err := d.decoder.Token(); err != io.EOF {
		return nil, newSyntaxError("invalid character after top-level value", d.decoder.InputOffset())
	}

	return item, nil
}

// This is not "legal", but it seems there is no other way to set the msg of syntax error.
func newSyntaxError(msg string, offset int64) *json.SyntaxError {
	err := &json.SyntaxError{
		Offset: offset,
	}

	msgField := reflect.ValueOf(err).Elem().Field(0 /* the msg field */)
	if msgField.Kind() == reflect.String {
		newMsgField := reflect.NewAt(msgField.Type(), unsafe.Pointer(msgField.UnsafeAddr())).Elem()
		newMsgField.SetString(msg)
	}

	return err
}

func (d *decoder) next() (any, error) {
	var token json.Token
	var err error

	if token, err = d.decoder.Token(); err != nil {
		return nil, err
	}

	return d.nextAfterToken(token)
}

func (d *decoder) nextAfterToken(token json.Token) (any, error) {
	var value any

	switch v := token.(type) {
	case bool, float64, json.Number, string, nil:
		value = v
	case json.Delim:
		switch v {
		case '{':
			{
				var object jsonObject[string, any]
				if d.opts.usePairs {
					object = NewPairs[string, any]()
				} else {
					m := NewMap[string, any]()
					m.SetDuplicatedKeyStrategy(d.opts.duplicatedKeyStrategy)
					object = m
				}
				if err := parseIntoObject[string, any](d, object, true); err != nil {
					return nil, err
				}
				value = object
			}
		case '[':
			{
				l := NewList[any]()
				if err := parseIntoArray[any](d, l); err != nil {
					return nil, err
				}
				value = l
			}
		}
	}

	return value, nil
}

// Array

type jsonArray[T any] interface {
	innerSlice() *[]T
}

func marshalArray[T any, A jsonArray[T]](array A) ([]byte, error) {
	slice := *array.innerSlice()
	if slice == nil {
		return []byte(`[]`), nil
	}

	var data bytes.Buffer
	enc := json.NewEncoder(&data)
	enc.SetEscapeHTML(false)
	err := enc.Encode(slice)
	return data.Bytes(), err
}

func parseIntoArray[T any, A jsonArray[T]](d *decoder, array A) error {
	*array.innerSlice() = nil

	for {
		token, err := d.decoder.Token()
		if err != nil {
			return err
		}

		// if meet ], the list parse ends
		delim, ok := token.(json.Delim)
		if ok && delim == ']' {
			return nil
		}

		var value T

		if v, err := d.nextAfterToken(token); err != nil {
			return err
		} else {
			value = v.(T)
		}

		*array.innerSlice() = append(*array.innerSlice(), value)
	}
}

func unmarshalArray[T any, A jsonArray[T]](data []byte, array A) error {
	if !isAny[T]() {
		return json.Unmarshal(data, array.innerSlice())
	}

	d := newDecoder(bytes.NewReader(data))

	token, err := d.decoder.Token()
	if err != nil {
		return err
	}

	if delim, ok := token.(json.Delim); !ok || delim != '[' {
		return &json.UnmarshalTypeError{
			Value: "non-array value",
			Type:  reflect.TypeOf(array).Elem(),
		}
	}

	return parseIntoArray[T](d, array)
}

// Object

type jsonObject[K comparable, V any] interface {
	GetByIndex(int) Pair[K, V]
	Add(K, V)
	Len() int
}

func marshalObject[K comparable, V any, O jsonObject[K, V]](object O) ([]byte, error) {
	if !isString[K]() {
		return nil, &json.UnsupportedTypeError{
			Type: reflect.TypeOf(object),
		}
	}

	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)

	buf.WriteByte('{')

	for i, length := 0, object.Len(); i < length; i++ {
		if i > 0 {
			buf.WriteByte(',')
		}

		pair := object.GetByIndex(i)

		// Key is string type, encoding will never fail
		enc.Encode(pair.Key)

		buf.WriteByte(':')

		if err := enc.Encode(pair.Value); err != nil {
			return nil, err
		}
	}

	buf.WriteByte('}')

	return buf.Bytes(), nil
}

func parseIntoObject[K comparable, V any, O jsonObject[K, V]](
	d *decoder, object O, valueIsAny bool,
) error {
	valueIsAny = valueIsAny || isAny[V]()

	for {
		token, err := d.decoder.Token()
		if err != nil {
			return err
		}

		// if meet }, the object parse ends
		delim, ok := token.(json.Delim)
		if ok && delim == '}' {
			return nil
		}

		// otherwise, we meet the key of a item
		key, _ := token.(string)

		var value V

		if valueIsAny { // if v is any, we parse it into our json value types
			if v, err := d.next(); err != nil {
				return err
			} else if v != nil {
				value = v.(V)
			}
		} else { // otherwise V is a real type, we can let std lib parsing it for us
			if err = d.decoder.Decode(&value); err != nil {
				return err
			}
		}

		var realKey K
		reflect.ValueOf(&realKey).Elem().SetString(key)

		object.Add(realKey, value)
	}
}

func unmarshalObject[K comparable, V any, O jsonObject[K, V]](
	data []byte, object O, option ...DecodeOption,
) error {
	if !isString[K]() {
		return &json.UnmarshalTypeError{
			Value: "any value",
			Type:  reflect.TypeOf(object).Elem(),
		}
	}

	d := newDecoder(bytes.NewReader(data), option...)

	token, err := d.decoder.Token()
	if err != nil {
		return err
	}

	if delim, ok := token.(json.Delim); !ok || delim != '{' {
		return &json.UnmarshalTypeError{
			Value: "non-object value",
			Type:  reflect.TypeOf(object).Elem(),
		}
	}

	return parseIntoObject[K, V](d, object, false)
}
