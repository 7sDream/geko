// Package geko provides GEneric Keep Order types.
//
// It's mainly used to solve the issue that in some scenarios, the field order
// in JSON object is meaningful, but when unmarshal into a normal map, this
// information will be lost. See [golang/go#27179].
//
// # Provided types
//
//   - [Map], and it's type alias [Object], to replace map.
//   - [Pairs], and it's type alias [ObjectItems], to replace map, when you need
//     to keep all values of duplicated key.
//   - [List], and it's type alias [Array] to replace slice.
//   - [Any] type, to replace the interface{}, it will use types above to
//     do json unmarshal.
//
// The [JSONUnmarshal] function is a shorthand for defined a [Any] and unmarshal
// data into it.
//
// # Example of JSON processing
//
//	result, _ := geko.JSONUnmarshal([]byte(`{"b": 1, "a": 2, "b": 3}`))
//	object := result.(geko.ObjectItems)
//	output, _ := json.Marshal(object)
//	fmt.Println(string(output)) // {"b":1,"a:2","b":3}
//
// If you do not want duplicated key in result, you can use [Pairs.ToMap], or
// use [UseObject] to let [JSONUnmarshal] do it for you:
//
//	result, _ := geko.JSONUnmarshal(
//		[]byte(`{"b": 1, "a": 2, "b": 3}`),
//		geko.UseObject(),
//	)
//	object, _ := result.(geko.Object)
//	object.Keys() // => ["b", "a"]
//	output, _ := json.Marshal(object)
//	fmt.Println(string(output)) // {"b":3,"a:2"}
//
// The [UseObject] option will make it use [Object] to unmarshal json object,
// instead of [ObjectItems]. [Object] will automatically deal with duplicated
// key for you. Maybe you think "b" should be 1, or "b" should appear after "a",
// This behavior can be adjusted by using [ObjectOnDuplicatedKey]
// with [DuplicatedKeyStrategy].
//
// [JSONUnmarshal] supports all json item, that's why it returns any. You can
// directly unmarshal into a [Object]/[ObjectItems] or [Array], if the type
// of input data is determined:
//
//	var arr geko.Array // or geko.Object for json object
//	_ := json.Unmarshal([]byte(`[1, 2, {"one": 1}, false]`), &arr)
//	object, _ := arr.Get(2).(geko.ObjectItems)
//	object.GetFirstOrZeroValue("one") // => 1
//
// But you can't customize [DecodeOptions] when doing this, it will always use
// default options.
//
// # Example for normally use
//
// Outside of JSON processing, these types can also be used simply as generic
// container types with insertion order preservation feature:
//
//	m := geko.NewMap[int, string]()
//	m.Set(1, "one")
//	m.Set(3, "three")
//	m.Set(2, "two")
//
//	m.Get(3) // "three", true
//	m.GetValueByIndex(1) // "three"
//
// There are many API for [Map], [List] and [Pairs], see their document for details.
//
// [golang/go#27179]: https://github.com/golang/go/issues/27179
package geko

import (
	"bytes"
	"encoding/json"
	"io"
	"reflect"
	"unsafe"
)

// DecodeOptions are options for controlling the behavior of [Any] unmarshaling.
//
// Zero value(default value) of it is:
//
//   - Do not use [json.Number] for json number, change it by apply [UseNumber]
//     option.
//   - Uses [ObjectItems] for json object, change it by apply [UseObject],
//     and change it back by apply [UseObjectItems] option.
//   - When [UseObject], the default [DuplicatedKeyStrategy] is
//     [UpdateValueKeepOrder], change it by apply
//     [ObjectOnDuplicatedKey] option.
//
// See also: [CreateDecodeOptions].
type DecodeOptions struct {
	useNumber             bool
	useObject             bool
	duplicatedKeyStrategy DuplicatedKeyStrategy
}

// DecodeOption is atom/modifier of [DecodeOptions].
type DecodeOption func(opts *DecodeOptions)

// CreateDecodeOptions creates a DecodeOptions by apply all option to the
// default decode option.
func CreateDecodeOptions(option ...DecodeOption) DecodeOptions {
	opts := DecodeOptions{}
	opts.Apply(option...)
	return opts
}

// Apply a option in current options.
func (opts *DecodeOptions) Apply(option ...DecodeOption) {
	for _, opt := range option {
		opt(opts)
	}
}

// UseNumber will change unmarshal behavior to using [json.Number] for json
// number.
func UseNumber(v bool) DecodeOption {
	return func(opts *DecodeOptions) {
		opts.useNumber = v
	}
}

// UseObject will change unmarshal behavior to using [Object] for JSON object.
func UseObject() DecodeOption {
	return func(opts *DecodeOptions) {
		opts.useObject = true
	}
}

// UseObject will change unmarshal behavior (back) to using [ObjectItem] for
// JSON object.
func UseObjectItem() DecodeOption {
	return func(opts *DecodeOptions) {
		opts.useObject = false
	}
}

// ObjectOnDuplicatedKey set the strategy when there are duplicated key in JSON
// object. Only effect when [UseObject] is applied.
//
// See document of [DuplicatedKeyStrategy] and its enum value for details.
func ObjectOnDuplicatedKey(strategy DuplicatedKeyStrategy) DecodeOption {
	return func(opts *DecodeOptions) {
		opts.duplicatedKeyStrategy = strategy
	}
}

type decoder struct {
	decoder *json.Decoder
	opts    DecodeOptions
}

func newDecoder(data []byte, opts DecodeOptions) *decoder {
	return &decoder{
		decoder: json.NewDecoder(bytes.NewReader(data)),
		opts:    opts,
	}
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
		return nil, newSyntaxError(
			"invalid character after top-level value",
			d.decoder.InputOffset(),
		)
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
				if d.opts.useObject {
					m := NewMap[string, any]()
					m.SetDuplicatedKeyStrategy(d.opts.duplicatedKeyStrategy)
					object = m
				} else {
					object = NewPairs[string, any]()
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
	// The behavior of the standard library is to clear the list
	// and we are consistent with it
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

func unmarshalArray[T any, A jsonArray[T]](data []byte, array A, option ...DecodeOption) error {
	if !isEmptyInterface[T]() {
		return json.Unmarshal(data, array.innerSlice())
	}

	d := newDecoder(data, CreateDecodeOptions(option...))

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
	// The behavior of the standard library is **do not** clear the map
	// and we are consistent with it.

	valueIsAny = valueIsAny || isEmptyInterface[V]()

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

	d := newDecoder(data, CreateDecodeOptions(option...))

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
