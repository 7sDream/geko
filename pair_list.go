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

func (kopl *PairList[K, V]) Append(pairs ...Pair[K, V]) {
	kopl.List = append(kopl.List, pairs...)
}

func (kopl *PairList[K, V]) Delete(key K) {
	kopl.Filter(func(p *Pair[K, V]) bool {
		return p.Key == key
	})
}

func (kopl *PairList[K, V]) DeleteByIndex(index int) {
	kopl.List = append(kopl.List[:index], kopl.List[index+1:]...)
}

func (kopl *PairList[K, V]) Clear() {
	kopl.List = nil
}

func (kopl *PairList[K, V]) Len() int {
	return len(kopl.List)
}

func (kopl *PairList[K, V]) Keys() []K {
	keys := make([]K, 0, kopl.Len())
	for i := 0; i < kopl.Len(); i++ {
		keys = append(keys, kopl.GetKeyByIndex(i))
	}
	return keys
}

func (kopl *PairList[K, V]) Values() []V {
	values := make([]V, 0, kopl.Len())
	for i := 0; i < kopl.Len(); i++ {
		values = append(values, kopl.GetValueByIndex(i))
	}
	return values
}

func (kopl *PairList[K, V]) ToMap(strategy DuplicateKeyStrategy) *Map[K, V] {
	kom := NewMap[K, V]()
	kom.SetDuplicateKeyStrategy(strategy)
	kom.Append(kopl.List...)
	return kom
}

func (kopl *PairList[K, V]) Dedup(strategy DuplicateKeyStrategy) {
	kopl.List = kopl.ToMap(strategy).Pairs().List
}

func (kopl *PairList[K, V]) Sort(lessFunc PairLessFunc[K, V]) {
	sort.SliceStable(kopl.List, func(i, j int) bool {
		return lessFunc(&kopl.List[i], &kopl.List[j])
	})
}

// Filter remove all item which make pred func return false.
//
// More efficient then `GetByIndex` + `DeleteByIndex` in a loop.
func (kopl *PairList[K, V]) Filter(pred PairFilterFunc[K, V]) {
	n := 0
	for i, length := 0, kopl.Len(); i < length; i++ {
		if pred(&kopl.List[i]) {
			kopl.List[n] = kopl.List[i]
			n++
		}
	}
	kopl.List = kopl.List[:n]
}

// MarshalJSON implements json.Marshaler interface.
// You should not call this directly, use json.Marshal(kom) instead.
func (kom PairList[K, V]) MarshalJSON() ([]byte, error) {
	return marshalObject[K, V](&kom)
}

// UnmarshalJSON implements json.Unmarshaler interface.
// You shouldn't call this directly, use json.Unmarshal(kom) instead.
func (kom *PairList[K, V]) UnmarshalJSON(data []byte) error {
	return unmarshalObject[K, V](data, kom, UsePairList(true))
}
