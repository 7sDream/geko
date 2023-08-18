package geko

import "sort"

// Pairs is a wrapper type of [][Pair][K, V].
//
// In JSON unmarshal, it will use the order of keys appear in JSON string,
// and marshal output will use the same order. But differ from [Map], it saves
// all items when their key is duplicated.
//
// When unmarshal from JSON into a *[Pairs][string, any], all JSON object will
// be stored in *[Pairs][string, any], all JSON array will be stored in
// *[List][any], instead of normal map[string]any and []any from std lib.
//
// Notice: Although this type behaves like a [Map], because it is only a slice
// internally, the performance of some APIs are not very good. It is best to
// keep this in mind when using it.
type Pairs[K comparable, V any] struct {
	List []Pair[K, V]
}

// ObjectItems is [Pairs] whose type parameters are specialized as
// [string, any], used to represent dynamic objects in JSON.
type ObjectItems = *Pairs[string, any]

// NewPairs creates a new empty list.
func NewPairs[K comparable, V any]() *Pairs[K, V] {
	return NewPairsFrom[K, V](nil)
}

// NewPairsWithCapacity likes [NewPairs], but init the inner container
// with a capacity to optimize memory allocate.
func NewPairsWithCapacity[K comparable, V any](capacity int) *Pairs[K, V] {
	return NewPairsFrom[K, V](make([]Pair[K, V], 0, capacity))
}

// NewPairsFrom create a List from a slice.
func NewPairsFrom[K comparable, V any](list []Pair[K, V]) *Pairs[K, V] {
	return &Pairs[K, V]{
		List: list,
	}
}

// Get values by key.
//
// Performance: O(n)
func (ps *Pairs[K, V]) Get(key K) []V {
	var values []V

	for i := range ps.List {
		p := &ps.List[i]
		if key == p.Key {
			values = append(values, p.Value)
		}
	}

	return values
}

// Has checks if a key exist in the list.
//
// Performance: O(n)
func (ps *Pairs[K, V]) Has(key K) bool {
	for i := range ps.List {
		if key == ps.List[i].Key {
			return true
		}
	}

	return false
}

// Count get appear times of a key.
//
// Performance: O(n)
func (ps *Pairs[K, V]) Count(key K) int {
	n := 0

	for i := range ps.List {
		if key == ps.List[i].Key {
			n++
		}
	}

	return n
}

// GetFirstOrZeroValue get first value by key, return a zero value of type V if
// key doesn't exist in list.
//
// Performance: O(n)
func (ps *Pairs[K, V]) GetFirstOrZeroValue(key K) (value V) {
	for i := range ps.List {
		p := &ps.List[i]
		if key == p.Key {
			value = p.Value
			break
		}
	}

	return
}

// GetFirstOrZeroValue get last value by key, return a zero value of type V if
// key doesn't exist in list.
//
// Performance: O(n)
func (ps *Pairs[K, V]) GetLastOrZeroValue(key K) (value V) {
	for i := ps.Len() - 1; i >= 0; i-- {
		p := &ps.List[i]
		if key == p.Key {
			value = p.Value
			break
		}
	}

	return
}

// GetKeyByIndex get key at index.
//
// You should make sure 0 <= i < Len(), panic if out of bound.
func (ps *Pairs[K, V]) GetKeyByIndex(index int) K {
	return ps.List[index].Key
}

// GetByIndex get key value pair at index.
//
// You should make sure 0 <= i < Len(), panic if out of bound.
func (ps *Pairs[K, V]) GetByIndex(index int) Pair[K, V] {
	return ps.List[index]
}

// GetValueByIndex get value at index.
//
// You should make sure 0 <= i < Len(), panic if out of bound.
func (ps *Pairs[K, V]) GetValueByIndex(index int) V {
	return ps.List[index].Value
}

// SetKeyByIndex changes key of item at index.
func (ps *Pairs[K, V]) SetKeyByIndex(index int, key K) {
	ps.List[index].Key = key
}

// SetValueByIndex changes value of item at index.
func (ps *Pairs[K, V]) SetValueByIndex(index int, value V) {
	ps.List[index].Value = value
}

// SetByIndex key and value at index.
func (ps *Pairs[K, V]) SetByIndex(index int, key K, value V) {
	ps.List[index] = CreatePair(key, value)
}

// Add a key value pair to the end of list.
func (ps *Pairs[K, V]) Add(key K, value V) {
	ps.List = append(ps.List, CreatePair(key, value))
}

// Append some key value pairs to the end of list.
func (ps *Pairs[K, V]) Append(pairs ...Pair[K, V]) {
	ps.List = append(ps.List, pairs...)
}

// Delete all item whose key is same as provided.
//
// Performance: O(n)
func (ps *Pairs[K, V]) Delete(key K) {
	ps.Filter(func(p *Pair[K, V]) bool {
		return p.Key != key
	})
}

// DeleteByIndex delete item at index.
//
// Performance: O(n)
func (ps *Pairs[K, V]) DeleteByIndex(index int) {
	ps.List = append(ps.List[:index], ps.List[index+1:]...)
}

// Clean this list.
func (ps *Pairs[K, V]) Clear() {
	ps.List = nil
}

// Len returns the size of list.
func (ps *Pairs[K, V]) Len() int {
	return len(ps.List)
}

// Keys returns all keys of the list.
//
// Performance: O(n).
func (ps *Pairs[K, V]) Keys() []K {
	keys := make([]K, 0, ps.Len())
	for i, length := 0, ps.Len(); i < length; i++ {
		keys = append(keys, ps.GetKeyByIndex(i))
	}
	return keys
}

// Values returns all values of the list.
//
// Performance: O(n).
func (ps *Pairs[K, V]) Values() []V {
	values := make([]V, 0, ps.Len())
	for i, length := 0, ps.Len(); i < length; i++ {
		values = append(values, ps.GetValueByIndex(i))
	}
	return values
}

// ToMap convert this list into a [Map], with provided [DuplicatedKeyStrategy].
func (ps *Pairs[K, V]) ToMap(strategy DuplicatedKeyStrategy) *Map[K, V] {
	m := NewMap[K, V]()
	m.SetDuplicatedKeyStrategy(strategy)
	m.Append(ps.List...)
	return m
}

// Dedup deduplicates this list by key.
//
// Implemented as converting it to a [Map] and back.
func (ps *Pairs[K, V]) Dedup(strategy DuplicatedKeyStrategy) {
	ps.List = ps.ToMap(strategy).Pairs().List
}

// Sort will reorder the list using the given less function.
func (ps *Pairs[K, V]) Sort(lessFunc PairLessFunc[K, V]) {
	sort.SliceStable(ps.List, func(i, j int) bool {
		return lessFunc(&ps.List[i], &ps.List[j])
	})
}

// Filter remove all item which make pred func return false.
//
// Performance: O(n). More efficient then [Pairs.GetByIndex] +
// [Pairs.DeleteByIndex] in a loop, which is O(n^2).
func (ps *Pairs[K, V]) Filter(pred PairFilterFunc[K, V]) {
	n := 0
	for i, length := 0, ps.Len(); i < length; i++ {
		if pred(&ps.List[i]) {
			ps.List[n] = ps.List[i]
			n++
		}
	}
	ps.List = ps.List[:n]
}

// MarshalJSON implements json.Marshaler interface.
// You should not call this directly, use json.Marshal(m) instead.
func (ps Pairs[K, V]) MarshalJSON() ([]byte, error) {
	return marshalObject[K, V](&ps)
}

// UnmarshalJSON implements json.Unmarshaler interface.
// You shouldn't call this directly, use json.Unmarshal(m) instead.
func (ps *Pairs[K, V]) UnmarshalJSON(data []byte) error {
	return unmarshalObject[K, V](data, ps, UseObjectItem())
}
