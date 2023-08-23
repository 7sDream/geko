package geko

import "encoding/json"

// Any is a wrapper for an any value. But when unmarshal, it uses our
// [Object]/[ObjectItems] and [Array] when meet JSON object and array.
//
// So the type of Any.Value after a [json.Unmarshal] can be:
// bool, float64/[json.Number], string, nil,
// [Object]/[ObjectItems], [Array].
//
// You can customize the unmarshal behavior by setting Any.Opts before call
// [json.Unmarshal].
//
// Notice: Usually you don't need to use this type directly. And, do not use
// this type on the value type parameter of the [Map], [Pairs] or [List].
// Because container types already handles standard any type specially,
// doing so will not only has no benefit, but also lose performance.
type Any struct {
	Value any
	Opts  DecodeOptions
}

// MarshalJSON implements [json.Marshaler] interface.
//
// You should not call this directly, use [json.Marshal] instead.
func (v Any) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.Value)
}

// UnmarshalJSON implements [json.Unmarshaler] interface.
//
// You shouldn't call this directly, use [json.Unmarshal]/[JSONUnmarshal]
// instead.
func (v *Any) UnmarshalJSON(data []byte) error {
	value, err := newDecoder(data, v.Opts).decode()
	if err == nil {
		v.Value = value
	}
	return err
}

// JSONUnmarshal is A convenience function for unmarshal JSON data into an
// [Any] and get the inner any value.
func JSONUnmarshal(data []byte, option ...DecodeOption) (any, error) {
	a := Any{Opts: CreateDecodeOptions(option...)}
	err := json.Unmarshal(data, &a)
	return a.Value, err
}
