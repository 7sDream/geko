# Geko

Geko is short for GEneric Keep Order, it provide generic map/list which keeps insert order of item.

When unmarshal from JSON string, the order is determined by their occurrence in the original string.

And of cause, when marshal back to JSON string, it will use the order too.

**Status**: WIP, feature not finished. API may not be the final version. Works in simple case, but not fully tested.

## Usage

### Unmarshal

```golang
// TODO
```

### Map

```golang
m := geko.NewMap[string, int]()
m.Set("one", 1)
m.Set("three", 2)
m.Set("two", 2)
m.Set("three", 3) // do not change order of key "three", it will stay ahead of "two".

for i, length := 0, m.Len(); i < length; i++ {
    pair := m.GetByIndex(i)
    fmt.Printf("%s: %d\n", pair.Key, pair.Value)
}

data, _ := json.Marshal(m)
fmt.Printf("marshal result: %s", string(data))
// Output:
// one: 1
// three: 3
// two: 2
// marshal result: {"one":1,"three":3,"two":2}
```

### List

```golang
// TODO
```

### PairList

```golang
// TODO
```

## TODO

- Tests.
- Package overview document and README.

## LICENSE

MIT. See LICENSE file.
