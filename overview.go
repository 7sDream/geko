// Package geko provides GEneric Keep Order types.
//
// It's mainly used to solve the issue that in some scenarios, the field order
// in JSON Object is meaningful, but when unmarshaling into a normal map, the
// order information will be erased. See [golang/go#27179].
//
// There are 3 types:
//
//   - [Map], and it's type alias [Object], to replace map.
//   - [Pairs], and it's type alias [ObjectItems], to replace map, when you need
//     to keep all values of duplicated key.
//   - [List], and it's type alias [Array] to replace slice.
//
// And a [JSONUnmarshal] function to replace [json.Unmarshal].
//
// # Example
//
//	result, _ := geko.JSONUnmarshal([]byte(`{"b": 1, "a": 2}`))
//	object, _ := result.(geko.Object)
//	object.Keys() // => ["b", "a"]
//	output, _ := json.Marshal(object)
//	fmt.Println(string(output)) // {"b":1,"a:2"}
//
// If you want use [Pairs] to deal with duplicated key:
//
//	result, _ := geko.JSONUnmarshal([]byte(`{"b": 1, "a": 2, "b": 3}`), geko.UsePairs(true))
//	object, _ := result.(geko.ObjectItems)
//	output, _ := json.Marshal(object)
//	fmt.Println(string(output)) // {"b":1,"a:2","b":3}
//
// You can use [OnDuplicatedKey] instead of [UsePairs], this will keep uses
// [Map], more specifically, [Object], and use [DuplicatedKeyStrategy] you
// provided to deal with duplicated key in object.
//
// Outside of JSON processing, these types can also be used simply as generic
// container types with insertion order preservation feature.
//
// [golang/go#27179]: https://github.com/golang/go/issues/27179
package geko
