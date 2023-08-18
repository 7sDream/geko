# geko

![license-badge] ![coverage-badge] [![document-badge]][document]

geko provides GEneric Keep Order types.

It's mainly used to solve the issue that in some scenarios, the field order in JSON object is meaningful, but when unmarshal into a normal map, these information will be lost. See [golang/go#27179].

There are many projects trying to solve it, but most of them lack some features that I need, see bellow.

## Features

- Not limited to JSON processing, can also be used as container with insertion order preservation feature.
- Generics, for better performance.
- Customizable strategy to deal with duplicated key, auto deduplication.
- Option to use json.Number to preserve the full precision of number field in unmarshaling.
- Fully tested, 100% coverage.

**Status**: Beta. All features I need are implemented and tested, But API design may not be the final version.

## Usage

### JSON processing

```golang
result, _ := geko.JSONUnmarshal([]byte(`{"a": 1, "b": 2, "a": 3}`))
object := result.(geko.ObjectItems)
output, _ := json.Marshal(object)
fmt.Println(string(output)) // {"a":1,"b":2,"a":3}
```

The `ObjectItems` is type alias of `*Pairs[string, any]`, which is a wrapper for `[]Pair[string, any]`. Because it's a slice under the hood, so it can store all items in JSON object.

If you don't want duplicated key in result, try `UseObject`:

```golang
result, _ := geko.JSONUnmarshal(
    []byte(`{"a": 1, "b": 2, "a": 3}`), 
    geko.UseObject(),
)
object := result.(geko.Object)
object.Keys() // => ["a", "b"]
output, _ := json.Marshal(object)
fmt.Println(string(output)) // {"a":3,"b":2}
```

`UseObject` will make `JSONUnmarshal` use `Object` to deal with JSON Object, it is alias of `*Map[string, any]`.

### Duplicated key strategy

You may find it weird that `a` has a value of `3`, this behavior can be controlled by add a option `geko.ObjectOnDuplicatedKey(strategy)`:

for input `{"a": 1, "b": 2, "a": 3}`

| strategy                 | result(space added) | note                      |
| :----------------------- | :------------------ | :------------------------ |
| `UpdateValueKeepOrder`   | `{"a": 3, "b": 2}`  | default strategy          |
| `UpdateValueUpdateOrder` | `{"b": 2, "a": 3}`  | keep the last occurrence  |
| `KeepValueUpdateOrder`   | `{"b": 2, "a": 1}`  |
| `Ignore`                 | `{"a": 1, "b": 2}`  | keep the first occurrence |

The `UpdateValueKeepOrder` is chosen as default strategy because it matches the behavior of NodeJS.

```text
> const obj = JSON.parse('{"a": 1, "b": 2, "a": 3}')
> obj.a
3
> JSON.stringify(obj)
'{"a":3,"b":2}'
> 
```

### Generic container

The type parameters do not limited to `[string, any]`, you can use other type you want.

#### Map

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
```

Outputs:

```text
one: 1
three: 3
two: 2
five: 5
four: 4
```

#### Pairs

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
```

Outputs:

```text
one: 1
three: 2
two: 2
three: 3
-----
one: 1
three: 2
two: 2
```

See [Document] for detail of all APIs.

## LICENSE

MIT. See LICENSE file.

[coverage-badge]: https://raw.githubusercontent.com/7sDream/geko/badges/.badges/master/coverage.svg
[document-badge]: https://img.shields.io/badge/-Document-blue?style=for-the-badge&logo=readthedocs
[license-badge]: https://img.shields.io/github/license/7sDream/geko?style=for-the-badge
[document]: https://pkg.go.dev/github.com/7sDream/geko
[golang/go#27179]: https://github.com/golang/go/issues/27179
