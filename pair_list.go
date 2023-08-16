package geko

import "sort"

// Wrapper type for [][Pair][K, V].
//
// In JSON unmarshal, it will use the order of keys appear in JSON string,
// and marshal output will use the same order. But differ from [Map], it saves
// all items when their key is duplicated.
//
// When unmarshal from JSON into a *[ParList][string, any], all JSON object will
// be stored in *ParList[string, any], all JSON array will be stored in
// *[List][any], instead of normal map[string]any and []any from std lib.
type PairList[K comparable, V any] struct {
	List []Pair[K, V]
}

// NewPairList creates a new empty list.
func NewPairList[K comparable, V any]() *PairList[K, V] {
	return NewPairListFrom[K, V](nil)
}

// NewPairListWithCapacity likes [NewPairList], but init the inner container
// with a capacity to optimize memory allocate.
func NewPairListWithCapacity[K comparable, V any](capacity int) *PairList[K, V] {
	return NewPairListFrom[K, V](make([]Pair[K, V], 0, capacity))
}

// NewPairListFrom create a List from a slice.
func NewPairListFrom[K comparable, V any](list []Pair[K, V]) *PairList[K, V] {
	return &PairList[K, V]{
		List: list,
	}
}

// Get values by key.
func (pl *PairList[K, V]) Get(key K) []V {
	var values []V

	for i := range pl.List {
		p := &pl.List[i]
		if key == p.Key {
			values = append(values, p.Value)
		}
	}

	return values
}

// Has checks if a key exist in the list.
func (pl *PairList[K, V]) Has(key K) bool {
	for i := range pl.List {
		if key == pl.List[i].Key {
			return true
		}
	}

	return false
}

// Count get appear times of a key.
func (pl *PairList[K, V]) Count(key K) int {
	n := 0

	for i := range pl.List {
		if key == pl.List[i].Key {
			n++
		}
	}

	return n
}

// GetFirstOrZeroValue get first value by key, return a zero value of type V if
// key doesn't exist in list.
func (pl *PairList[K, V]) GetFirstOrZeroValue(key K) (value V) {
	for i := range pl.List {
		p := &pl.List[i]
		if key == p.Key {
			value = p.Value
			break
		}
	}

	return
}

// GetFirstOrZeroValue get last value by key, return a zero value of type V if
// key doesn't exist in list.
func (pl *PairList[K, V]) GetLastOrZeroValue(key K) (value V) {
	for i := pl.Len() - 1; i >= 0; i-- {
		p := &pl.List[i]
		if key == p.Key {
			value = p.Value
			break
		}
	}

	return
}

func (pl *PairList[K, V]) GetKeyByIndex(index int) K {
	return pl.List[index].Key
}

func (pl *PairList[K, V]) GetByIndex(index int) Pair[K, V] {
	return pl.List[index]
}

func (pl *PairList[K, V]) GetValueByIndex(index int) V {
	return pl.List[index].Value
}

func (pl *PairList[K, V]) Add(key K, value V) {
	pl.List = append(pl.List, Pair[K, V]{key, value})
}

func (pl *PairList[K, V]) Append(pairs ...Pair[K, V]) {
	pl.List = append(pl.List, pairs...)
}

func (pl *PairList[K, V]) Delete(key K) {
	pl.Filter(func(p *Pair[K, V]) bool {
		return p.Key == key
	})
}

func (pl *PairList[K, V]) DeleteByIndex(index int) {
	pl.List = append(pl.List[:index], pl.List[index+1:]...)
}

func (pl *PairList[K, V]) Clear() {
	pl.List = nil
}

func (pl *PairList[K, V]) Len() int {
	return len(pl.List)
}

func (pl *PairList[K, V]) Keys() []K {
	keys := make([]K, 0, pl.Len())
	for i := 0; i < pl.Len(); i++ {
		keys = append(keys, pl.GetKeyByIndex(i))
	}
	return keys
}

func (pl *PairList[K, V]) Values() []V {
	values := make([]V, 0, pl.Len())
	for i := 0; i < pl.Len(); i++ {
		values = append(values, pl.GetValueByIndex(i))
	}
	return values
}

func (pl *PairList[K, V]) ToMap(strategy DuplicatedKeyStrategy) *Map[K, V] {
	m := NewMap[K, V]()
	m.SetDuplicatedKeyStrategy(strategy)
	m.Append(pl.List...)
	return m
}

func (pl *PairList[K, V]) Dedup(strategy DuplicatedKeyStrategy) {
	pl.List = pl.ToMap(strategy).Pairs().List
}

func (pl *PairList[K, V]) Sort(lessFunc PairLessFunc[K, V]) {
	sort.SliceStable(pl.List, func(i, j int) bool {
		return lessFunc(&pl.List[i], &pl.List[j])
	})
}

// Filter remove all item which make pred func return false.
//
// More efficient then `GetByIndex` + `DeleteByIndex` in a loop.
func (pl *PairList[K, V]) Filter(pred PairFilterFunc[K, V]) {
	n := 0
	for i, length := 0, pl.Len(); i < length; i++ {
		if pred(&pl.List[i]) {
			pl.List[n] = pl.List[i]
			n++
		}
	}
	pl.List = pl.List[:n]
}

// MarshalJSON implements json.Marshaler interface.
// You should not call this directly, use json.Marshal(m) instead.
func (m PairList[K, V]) MarshalJSON() ([]byte, error) {
	return marshalObject[K, V](&m)
}

// UnmarshalJSON implements json.Unmarshaler interface.
// You shouldn't call this directly, use json.Unmarshal(m) instead.
func (m *PairList[K, V]) UnmarshalJSON(data []byte) error {
	return unmarshalObject[K, V](data, m, UsePairList(true))
}
