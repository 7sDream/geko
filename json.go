package geko

import (
	"bytes"
	"encoding/json"
	"io"
	"reflect"
	"unsafe"
)

// Decode likes json.Unmarshal, but uses our Map and List when meet JSON object and array.
//
// So the returned value can be:
// bool, float64/json.Number, string, nil, Map[string]any/PairList[string]any, List[any].
//
// The `any` value in the above container can only be the above type, recursive.
func JSONUnmarshal(data []byte, option ...Option) (any, error) {
	return newDecoder(bytes.NewReader(data), option...).decode()
}

type options struct {
	useNumber   bool
	usePairList bool
}

type Option func(opts *options)

func UseNumber(v bool) Option {
	return func(opts *options) {
		opts.useNumber = v
	}
}

func UsePairList(v bool) Option {
	return func(opts *options) {
		opts.usePairList = v
	}
}

type decoder struct {
	decoder *json.Decoder
	opts    options
}

func newDecoder(r io.Reader, option ...Option) *decoder {
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

	item, err := d.next(nil)
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

func (d *decoder) next(escapeHTML *bool) (any, error) {
	var token json.Token
	var err error

	if token, err = d.decoder.Token(); err != nil {
		return nil, err
	}

	return d.nextByToken(token, escapeHTML)
}

func (d *decoder) nextByToken(token json.Token, escapeHTML *bool) (any, error) {
	var value any

	switch v := token.(type) {
	case bool, float64, json.Number, string, nil:
		value = v
	case json.Delim:
		switch v {
		case '{':
			{
				var object jsonObjectLike[string, any]
				if d.opts.usePairList {
					object = NewPairList[string, any]()
				} else {
					object = NewMap[string, any]()
				}
				if escapeHTML != nil {
					object.SetEscapeHTML(*escapeHTML)
				}
				if err := parseIntoObject(d, object, true); err != nil {
					return nil, err
				}
				value = object
			}
		case '[':
			{
				kol := NewList[any]()
				if escapeHTML != nil {
					kol.SetEscapeHTML(*escapeHTML)
				}
				if err := parseIntoArray(d, kol); err != nil {
					return nil, err
				}
				value = kol
			}
		}
	}

	return value, nil
}

// Array

type jsonArrayLike[T any] interface {
	EscapeHTML() bool
	innerSlice() *[]T
}

func marshalArray[T any, A jsonArrayLike[T]](array A) ([]byte, error) {
	var data bytes.Buffer
	encoder := json.NewEncoder(&data)
	encoder.SetEscapeHTML(array.EscapeHTML())

	err := encoder.Encode(*array.innerSlice())
	return data.Bytes(), err
}

func parseIntoArray[T any, A jsonArrayLike[T]](d *decoder, array A) error {
	escape := array.EscapeHTML()

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

		if v, err := d.nextByToken(token, &escape); err != nil {
			return err
		} else {
			value = v.(T)
		}

		*array.innerSlice() = append(*array.innerSlice(), value)
	}
}

func unmarshalArray[T any, A jsonArrayLike[T]](data []byte, array A) error {
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

	return parseIntoArray(d, array)
}

// Object

type jsonObjectLike[K comparable, V any] interface {
	EscapeHTML() bool
	SetEscapeHTML(bool)
	GetByIndex(int) Pair[K, V]
	Set(K, V)
	Len() int
}

func marshalObject[K comparable, V any, O jsonObjectLike[K, V]](object O) ([]byte, error) {
	if !underlyingIsString[K]() {
		return nil, &json.UnsupportedTypeError{
			Type: reflect.TypeOf(object),
		}
	}

	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(object.EscapeHTML())

	buf.WriteByte('{')

	for i := 0; i < object.Len(); i++ {
		if i > 0 {
			buf.WriteByte(',')
		}

		pair := object.GetByIndex(i)

		if err := enc.Encode(pair.Key); err != nil {
			return nil, err
		}

		buf.WriteByte(':')

		if err := enc.Encode(pair.Value); err != nil {
			return nil, err
		}
	}

	buf.WriteByte('}')

	return buf.Bytes(), nil
}

func parseIntoObject[K comparable, V any, O jsonObjectLike[K, V]](
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
			escape := object.EscapeHTML()
			if v, err := d.next(&escape); err != nil {
				return err
			} else {
				value = v.(V)
			}
		} else { // otherwise V is a real type, we can let std lib parsing it for us
			if err = d.decoder.Decode(&value); err != nil {
				return err
			}
		}

		var realKey K
		reflect.ValueOf(&realKey).Elem().SetString(key)

		object.Set(realKey, value)
	}
}

func unmarshalObject[K comparable, V any, O jsonObjectLike[K, V]](data []byte, object O) error {
	if !underlyingIsString[K]() {
		return &json.UnmarshalTypeError{
			Type: reflect.TypeOf(object).Elem(),
		}
	}

	d := newDecoder(bytes.NewReader(data))

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

	return parseIntoObject(d, object, false)
}
