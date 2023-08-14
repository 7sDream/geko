package geko

// Map is a map, its kv pair keeps insert order.
// In JSON Unmarshal and marshal, it will use the order of keys appear in JSON string and output as is.
//
// When Unmarshal from json into a `Map[string]any`, all json object will store in `Map[string]any`,
// all json array will store in `List[any]`, instead of normal `map[string]any` and `[]any` from stdlib.
type Map[K comparable, V any] struct {
	order []K
	inner map[K]V

	escapeHTML bool
}

// NewMap creates a new empty ordered map
func NewMap[K comparable, V any]() *Map[K, V] {
	return &Map[K, V]{
		escapeHTML: true, // keep sync with golang std lib
	}
}

// NewMapWithCapacity likes NewMap, but init the inner container with a capacity to optimize memory allocate.
func NewMapWithCapacity[K comparable, V any](capacity int) *Map[K, V] {
	kom := NewMap[K, V]()
	kom.order = make([]K, 0, capacity)
	kom.inner = make(map[K]V, capacity)
	return kom
}

// Get a value by key. the second return value is true if the key exists, otherwise false.
func (kom *Map[K, V]) Get(key K) (V, bool) {
	v, exist := kom.inner[key]
	return v, exist
}

// GetOrZeroValue return stored value by key, or te zero value of value type if key not exist.
func (kom *Map[K, V]) GetOrZeroValue(key K) V {
	return kom.inner[key]
}

// GetKeyByIndex get key by it's index in order.
// You should make sure 0 <= i < Len(), panic if out of bound.
func (kom *Map[K, V]) GetKeyByIndex(index int) K {
	return kom.order[index]
}

// GetByIndex get the key and value by index of key order.
// You should make sure 0 <= i < Len(), panic if out of bound.
func (kom *Map[K, V]) GetByIndex(index int) Pair[K, V] {
	k := kom.GetKeyByIndex(index)
	return Pair[K, V]{Key: k, Value: kom.GetOrZeroValue(k)}
}

// GetValueByIndex get the value by index of key order.
// You should make sure 0 <= i < Len(), panic if out of bound.
func (kom *Map[K, V]) GetValueByIndex(index int) V {
	k := kom.GetKeyByIndex(index)
	return kom.GetOrZeroValue(k)
}

// Set a value by key.
// Called with an already exist key will not change it's order.
// If you want move it to the end, call Delete before Set.
func (kom *Map[K, V]) Set(key K, value V) {
	if kom.inner == nil {
		kom.inner = make(map[K]V)
	}

	if _, exist := kom.inner[key]; !exist {
		kom.order = append(kom.order, key)
	}

	kom.inner[key] = value
}

func (kom *Map[K, V]) Extend(pairs ...Pair[K, V]) {
	for _, pair := range pairs {
		kom.Set(pair.Key, pair.Value)
	}
}

// Delete a item by key.
func (kom *Map[K, V]) Delete(key K) {
	_, exist := kom.inner[key]
	if !exist {
		return
	}

	for i, k := range kom.order {
		if k == key {
			kom.DeleteByIndex(i)
			return
		}
	}
}

// Delete a item by it's index in order.
// You should make sure 0 <= i < Len(), panic if out of bound.
func (kom *Map[K, V]) DeleteByIndex(index int) {
	key := kom.order[index]
	kom.order = append(kom.order[:index], kom.order[index+1:]...)
	delete(kom.inner, key)
}

// Clean this map.
func (kom *Map[K, V]) Clear() {
	kom.order = nil
	kom.inner = nil
}

// Len returns the size of map.
func (kom *Map[K, V]) Len() int {
	return len(kom.inner)
}

// Keys returns the keys of ordered map, in current order.
// This will copy all keys, so you can modify it if you wish.
// If you want iterate over the map, maybe Len + GetByIndex is a better choice.
func (kom *Map[K, V]) Keys() []K {
	// copy to avoid user modify the order.
	keys := make([]K, 0, kom.Len())
	copy(keys, kom.order)
	return keys
}

func (kom *Map[K, V]) Values() []V {
	values := make([]V, 0, kom.Len())
	for i := 0; i < kom.Len(); i++ {
		values = append(values, kom.GetValueByIndex(i))
	}
	return values
}

// ToPairs gives you all data the map stored as a list of pair, in current order.
// Do not use this to iterate the map, use GetByIndex instead.
func (kom *Map[K, V]) Pairs() *PairList[K, V] {
	pairs := NewPairListWithCapacity[K, V](kom.Len())

	for i := 0; i < kom.Len(); i++ {
		pairs.List = append(pairs.List, kom.GetByIndex(i))
	}

	return pairs
}

// Sort will reorder the map using the given less function.
func (kom *Map[K, V]) Sort(lessFunc PairLessFunc[K, V]) {
	pairs := kom.Pairs()

	pairs.Sort(lessFunc)

	for i := 0; i < kom.Len(); i++ {
		kom.order[i] = pairs.List[i].Key
	}
}

// MarshalJSON implements json.Marshaler interface.
// You should not call this directly, use json.Marshal(kom) instead.
func (kom Map[K, V]) MarshalJSON() ([]byte, error) {
	return marshalObject[K, V](&kom)
}

// UnmarshalJSON implements json.Unmarshaler interface.
// You shouldn't call this directly, use json.Unmarshal(kom) instead.
func (kom *Map[K, V]) UnmarshalJSON(data []byte) error {
	return unmarshalObject[K, V](data, kom)
}
