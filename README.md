# Geko

Geko is short for GEneric Keep Order, it provide generic map/list which keeps insert order of item.

When unmarshal from JSON string, the order is determined by their occurrence in the original string.

And of cause, when marshal back to JSON string, it will use the order too.

**Status**: WIP, very early version. API may not be the final version. Works in simple case, but not fully tested.

## Usage

### Unmarshal and marshal

```golang
// TODO
```

### Map

```golang
import (
    "fmt"

    "github.com/7sDream/geko"
)

func main() {
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
}
```

### List

```golang
// TODO
```

### PairList

```golang
// TODO
```

### Decoder

```golang
// TODO
```

## Known Issue

Due to the design of `json.Marshaler` interface is `MarshalJSON() ([]byte, error)`, a custom type have no way to get options set into `json.Encoder`, like `EscapeHTML`, when marshal itself.

So Map/List/PairList type will (actually, can only) ignore the options in `json.Encoder`, but uses a global option defined in this module instead. I recommend you set this option at begin or init function of your project as your needs.  

See `geko.EscapeHTML()` and `geko.SetEscapeHTML()` for detail.

## TODO

- Doc/Comments.
- Tests.

## LICENSE

MIT. See LICENSE file.
