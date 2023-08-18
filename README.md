# Geko

[Document]

geko provides GEneric Keep Order map types.

It's mainly used to solve the issue that in some scenarios, the field order in JSON object is meaningful, but when unmarshal into a normal map, the order information will be erased. See [golang/go#27179].

Outside of JSON processing, the map can also be used simply as generic container types with insertion order preservation feature.

**Status**: Beta. All feature I need is implemented and tested. But API design may not be the final version.

## Usage

### JSON processing

```golang
result, _ := geko.JSONUnmarshal([]byte(`{"a": 1, "b": 2, "a": 3}`))
object := result.(geko.ObjectItems)
output, _ := json.Marshal(object)
fmt.Println(string(output)) // {"a":1,"b":2,"a":3}
```

The `ObjectItems` is type alias of `*Pairs[string, any]`, which is a wrapper for `[]Pair[string, any]`. Because there is a slice under the hood, so it can store all items in JSON object.

If you do not want duplicated key in result, try `UseObject`:

```golang
result, _ := geko.JSONUnmarshal([]byte(`{"a": 1, "b": 2, "a": 3}`), geko.UseObject())
object := result.(geko.Object)
object.Keys() // => ["a", "b"]
output, _ := json.Marshal(object)
fmt.Println(string(output)) // {"a":3,"b":2}
```

### Duplicated key strategy

You may find it weird that `a` has a value of `3`, this behavior can be controlled by add a option `geko.ObjectOnDuplicatedKey(strategy)`:

for input `{"a": 1, "b": 2, "a": 3}`

strategy | result(space added) | note
:------- | :----- | :----
UpdateValueKeepOrder | `{"a": 3, "b": 2}` | default strategy
UpdateValueUpdateOrder | `{"b": 2, "a": 3}` | use last
KeepValueUpdateOrder | `{"b": 2, "a": 1}` |
Ignore | `{"a": 1, "b": 2}` | use first

The `UpdateValueKeepOrder` is chosen as default strategy because it matches the behavior of NodeJS.

```text
> const obj = JSON.parse('{"a": 1, "b": 2, "a": 3}')
> obj.a
3
> JSON.stringify(obj)
'{"a":3,"b":2}'
> 
```

### Use container type directly

Those type can also be used directly:

```golang
m := geko.NewMap[string, int]()

m.Set("one", 1)
m.Set("three", 2)
m.Set("two", 2)
m.Set("three", 3) // Set always do not change order of existed key, so "three" will stay ahead of "two".
m.Set("four", 0)
m.Set("five", 5)

m.SetDuplicatedKeyStrategy(geko.UpdateValueUpdateOrder)
m.Add("four", 4) // Add will follow DuplicatedKeyStrategy, so now four is last key, and it's value is 4

for i, length := 0, m.Len(); i < length; i++ {
    pair := m.GetByIndex(i)
    fmt.Printf("%s: %d\n", pair.Key, pair.Value)
}

// Output:
// one: 1
// three: 3
// two: 2
// five: 5
// four: 4
```

``` golang
m := geko.NewPairs[string, int]()

m.Add("one", 1)
m.Add("three", 2)
m.Add("two", 2)
m.Add("three", 3)
for i, length := 0, m.Len(); i < length; i++ {
    pair := m.GetByIndex(i)
    fmt.Printf("%s: %d\n", pair.Key, pair.Value)
}

fmt.Println("-----")

m.Dedup(geko.Ignore)
for i, length := 0, m.Len(); i < length; i++ {
    pair := m.GetByIndex(i)
    fmt.Printf("%s: %d\n", pair.Key, pair.Value)
}

// Output:
// one: 1
// three: 2
// two: 2
// three: 3
// -----
// one: 1
// three: 2
// two: 2
```

There are many APIs for Map and Pairs, like Sort, Filter or Dedup. See [Document] for details.

## LICENSE

MIT. See LICENSE file.

[Document]: https://pkg.go.dev/github.com/7sDream/geko
[golang/go#27179]: https://github.com/golang/go/issues/27179
