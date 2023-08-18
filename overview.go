// Package geko provides some common order-preserving types, mainly to solve the
// problem that in some scenarios, the field order in JSON Object is meaningful,
// but when unmarshaling into a a normal map, the order information will be
// erased.
//
// There are 3 types:
// - [Map] to replace map.
// - [Pairs] to replace map, when you need to keep all values of duplicated key
// - [List] to replace list.
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
package geko
