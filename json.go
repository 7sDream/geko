package geko

import (
	"bytes"
	"encoding/json"
	"io"
	"reflect"
	"unsafe"
)

// ===== Encoder =====

// JSONUnmarshal likes json.Unmarshal, but uses our Map and List when meet JSON object and array.
//
// So the returned value can be:
// bool, float64/json.Number, string, nil, Map[string]any/PairList[string]any, List[any].
//
// The `any` value in the above container can only be the above type, recursive.
func JSONMarshal(v any, option ...EncodeOption) ([]byte, error) {
	return nil, newEncoder(option...).encode(v)
}

type encoderOptions struct {
	escapeHTML           bool
	indentPrefix, indent string
}

type EncodeOption func(opts *encoderOptions)

func EscapeHTML(v bool) EncodeOption {
	return func(opts *encoderOptions) {
		opts.escapeHTML = v
	}
}

func WithIndent(prefix string, indent string) EncodeOption {
	return func(opts *encoderOptions) {
		opts.indentPrefix = prefix
		opts.indent = indent
	}
}

type encoder struct {
	buf     bytes.Buffer
	encoder *json.Encoder
	opts    encoderOptions
}

func newEncoder(option ...EncodeOption) *encoder {
	encoder := &encoder{}

	for _, opt := range option {
		opt(&encoder.opts)
	}

	return encoder
}

type jsonMarshalerExt interface {
	jsonMarshalExt(e *encoder) ([]byte, error)
}

func (e *encoder) encode(v any) error {
	e.encoder = json.NewEncoder(&e.buf)
	e.encoder.SetEscapeHTML(e.opts.escapeHTML)
	e.encoder.SetIndent(e.opts.indentPrefix, e.opts.indent)

	if ext, ok := v.(jsonMarshalerExt); ok {
		return ext.jsonMarshalExt(e)
	}

	switch v.(type) {
	case bool, float64, json.Number, string, nil:
		{
			return e.encoder.Encode(v)
		}
	case map[string]any:
		{
			// TODO
		}
	case []any:
		{
			// TODO
		}
	}

	return nil
}

func encodeObject[K comparable, V any, O jsonObjectLike[K, V]](
	e *encoder, object jsonObjectLike[string, any],
) error {
	if !IsString[K]() {
		return &json.UnsupportedTypeError{
			Type: reflect.TypeOf(object),
		}
	}

	e.buf.WriteByte('{')

	for i := 0; i < object.Len(); i++ {
		if i > 0 {
			e.buf.WriteByte(',')
		}

		pair := object.GetByIndex(i)

		if err := e.encode(pair.Key); err != nil {
			return err
		}

		e.buf.WriteByte(':')

		if err := e.encode(pair.Value); err != nil {
			return err
		}
	}

	e.buf.WriteByte('}')

	return nil
}

// ===== Decoder =====

// JSONUnmarshal likes json.Unmarshal, but uses our Map and List when meet JSON object and array.
//
// So the returned value can be:
// bool, float64/json.Number, string, nil, Map[string]any/PairList[string]any, List[any].
//
// The `any` value in the above container can only be the above type, recursive.
func JSONUnmarshal(data []byte, option ...DecodeOption) (any, error) {
	return newDecoder(bytes.NewReader(data), option...).decode()
}

type decodeOptions struct {
	useNumber   bool
	usePairList bool
}

type DecodeOption func(opts *decodeOptions)

func UseNumber(v bool) DecodeOption {
	return func(opts *decodeOptions) {
		opts.useNumber = v
	}
}

func UsePairList(v bool) DecodeOption {
	return func(opts *decodeOptions) {
		opts.usePairList = v
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
				var object jsonObjectLike[string, any]
				if d.opts.usePairList {
					object = NewPairList[string, any]()
				} else {
					object = NewMap[string, any]()
				}
				if err := parseIntoObject[string, any](d, object, true); err != nil {
					return nil, err
				}
				value = object
			}
		case '[':
			{
				kol := NewList[any]()
				if err := parseIntoArray[any](d, kol); err != nil {
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
	innerSlice() *[]T
}

func marshalArray[T any, A jsonArrayLike[T]](array A) ([]byte, error) {
	var data bytes.Buffer
	encoder := json.NewEncoder(&data)

	err := encoder.Encode(*array.innerSlice())
	return data.Bytes(), err
}

func parseIntoArray[T any, A jsonArrayLike[T]](d *decoder, array A) error {
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

	return parseIntoArray[T](d, array)
}

// Object

type jsonObjectLike[K comparable, V any] interface {
	GetByIndex(int) Pair[K, V]
	Set(K, V)
	Len() int
}

func marshalObject[K comparable, V any, O jsonObjectLike[K, V]](object O) ([]byte, error) {
	if !IsString[K]() {
		return nil, &json.UnsupportedTypeError{
			Type: reflect.TypeOf(object),
		}
	}

	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)

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
			if v, err := d.next(); err != nil {
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

func unmarshalObject[K comparable, V any, O jsonObjectLike[K, V]](data []byte, object O, option ...DecodeOption) error {
	if !IsString[K]() {
		return &json.UnmarshalTypeError{
			Type: reflect.TypeOf(object).Elem(),
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
