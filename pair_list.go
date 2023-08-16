package geko

import "sort"

// Wrapper type for []Pair[K, V], it can be unmarshal/marshal from/to a JSON
// **object**. It will use the order of keys appear in JSON string and output
// as is. It saves all items even if their key is duplicated.
//
// When Unmarshal from json into a `ParList[string, any]`, all inner JSON object
// will be stored in `*ParList[string, any]`, all JSON array will be stored in
// `*List[any]`, instead of normal `map[string]any` and `[]any` from std lib.
type PairList[K comparable, V any] struct {
	List []Pair[K, V]
}

func NewPairList[K comparable, V any]() *PairList[K, V] {
	return NewPairListFrom[K, V](nil)
}

func NewPairListWithCapacity[K comparable, V any](capacity int) *PairList[K, V] {
	return NewPairListFrom[K, V](make([]Pair[K, V], 0, capacity))
}

func NewPairListFrom[K comparable, V any](list []Pair[K, V]) *PairList[K, V] {
	return &PairList[K, V]{
		List: list,
	}
}

func (kopl *PairList[K, V]) Get(key K) []V {
	var values []V

	for _, pair := range kopl.List {
		if key == pair.Key {
			values = append(values, pair.Value)
		}
	}

	return values
}

func (kopl *PairList[K, V]) GetFirstOrZeroValue(key K) (value V) {
	for _, pair := range kopl.List {
		if key == pair.Key {
			value = pair.Value
		}
	}
	return
}

func (kopl *PairList[K, V]) GetKeyByIndex(index int) K {
	return kopl.List[index].Key
}

func (kopl *PairList[K, V]) GetByIndex(index int) Pair[K, V] {
	return kopl.List[index]
}

func (kopl *PairList[K, V]) GetValueByIndex(index int) V {
	return kopl.List[index].Value
}

func (kopl *PairList[K, V]) Set(key K, value V) {
	kopl.List = append(kopl.List, Pair[K, V]{key, value})
}

func (kopl *PairList[K, V]) Extend(pairs ...Pair[K, V]) {
	kopl.List = append(kopl.List, pairs...)
}

func (kopl *PairList[K, V]) Delete(key K) {
	for i := kopl.Len() - 1; i > 0; i-- {
		if key == kopl.GetKeyByIndex(i) {
			kopl.DeleteByIndex(i)
		}
	}
}

func (kopl *PairList[K, V]) DeleteByIndex(index int) {
	kopl.List = append(kopl.List[:index], kopl.List[index+1:]...)
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
