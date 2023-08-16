# Geko

Geko is short for GEneric Keep Order, it provide generic map/list which keeps insert order of item.

When unmarshal from JSON string, the order is determined by their occurrence in the original string.

And of cause, when marshal back to JSON string, it will use the order too.

**Status**: WIP, very early version. API may not be the final version. Works in simple case, but not fully tested.

## Usage

### Unmarshal

```golang
// TODO
```

### Map

```golang
kom := geko.NewMap[string, int]()
kom.Set("one", 1)
kom.Set("three", 2)
kom.Set("two", 2)
kom.Set("three", 3) // do not change order of key "three", it will stay ahead of "two".

for i := 0; i < kom.Len(); i++ {
    pair := kom.GetByIndex(i)
    fmt.Printf("%s: %d\n", pair.Key, pair.Value)
}

data, _ := json.Marshal(kom)
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

- Doc/Comments.
- Tests.

## LICENSE

MIT. See LICENSE file.
