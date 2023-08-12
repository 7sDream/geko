package geko_test

import (
	"encoding/json"
	"fmt"

	"github.com/7sDream/geko"
)

// Use is to iterate over the map.
func ExampleMap_GetByIndex() {
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
