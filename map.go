package geko

// DuplicateKeyStrategy controls the behavior of [Map.Set] when meet a
// duplicate key. Default strategy is [UpdateValueKeepOrder].
//
// If you want store all values of duplicate key, use [PairList] type instead.
type DuplicateKeyStrategy uint8

const (
	// UpdateValueKeepOrder will use new value, but keep the key order.
	//
	// {"a": 1, "b": 2, "a": 3} => {"a": 3, "b": 2}
	//
	// This is the default strategy.
	UpdateValueKeepOrder DuplicateKeyStrategy = iota
	// UpdateValueUpdateOrder will use new value, and move the key to last.
	//
	// {"a": 1, "b": 2, "a": 3} => {"b": 2, "a": 3}
	UpdateValueUpdateOrder
	// KeepValueUpdateOrder will keep the value not change, but move the key to
	// last.
	//
	// {"a": 1, "b": 2, "a": 3} => {"b": 2, "a": 1}
	KeepValueUpdateOrder
	// Ignore will ignore the wholes set option when key duplicated.
	//
	// {"a": 1, "b": 2, "a": 3} => {"a": 1, "b": 2}
	Ignore
)

// Map is a map, in which the kv pairs will keep order of their insert.
//
// In JSON unmarshal, it will use the order of keys appear in JSON string,
// and marshal output will use the same order.
//
// When unmarshal from JSON into a *[Map][string, any], all JSON object will be
// stored in *[Map][string, any], all JSON array will be stored in *[List][any],
// instead of normal map[string]any and []any from std lib.
//
// You can [Map.SetDuplicateKeyStrategy] before call [json.Unmarshal] to control
// the behavior when object has duplicate key in your JSON string data.
//
// If you do not sure the outmost item is object, see [JSONUnmarshal] function.
type Map[K comparable, V any] struct {
	order []K
	inner map[K]V

	onDuplicateKey DuplicateKeyStrategy
}

// NewMap creates a new empty map.
func NewMap[K comparable, V any]() *Map[K, V] {
	return &Map[K, V]{}
}

// NewMapWithCapacity likes [NewMap], but init the inner container with a
// capacity to optimize memory allocate.
func NewMapWithCapacity[K comparable, V any](capacity int) *Map[K, V] {
	m := NewMap[K, V]()
	m.order = make([]K, 0, capacity)
	m.inner = make(map[K]V, capacity)
	return m
}

// DuplicateKeyStrategy get current strategy when Set with a duplicate key.
//
// See document of [DuplicateKeyStrategy] and its enum value for detail.
func (m *Map[K, V]) DuplicateKeyStrategy() DuplicateKeyStrategy {
	return m.onDuplicateKey
}

// SetDuplicateKeyStrategy set strategy when [Map.Set] with a duplicate key.
//
// See document of [DuplicateKeyStrategy] and its enum value for detail.
func (m *Map[K, V]) SetDuplicateKeyStrategy(strategy DuplicateKeyStrategy) {
	m.onDuplicateKey = strategy
}

// Get a value by key. The second return value is true if the key exists,
// otherwise false.
func (m *Map[K, V]) Get(key K) (V, bool) {
	v, exist := m.inner[key]
	return v, exist
}

// Has checks if key is in the map.
func (m *Map[K, V]) Has(key K) bool {
	_, exist := m.inner[key]
	return exist
}

// GetOrZeroValue return stored value by key, or the zero value of value type
// if key not exist.
func (m *Map[K, V]) GetOrZeroValue(key K) V {
	return m.inner[key]
}

// GetKeyByIndex get key by it's index in order.
//
// You should make sure 0 <= i < Len(), panic if out of bound.
func (m *Map[K, V]) GetKeyByIndex(index int) K {
	return m.order[index]
}

// GetByIndex get the key and value by index of key order.
//
// You should make sure 0 <= i < Len(), panic if out of bound.
func (m *Map[K, V]) GetByIndex(index int) Pair[K, V] {
	k := m.GetKeyByIndex(index)
	return Pair[K, V]{Key: k, Value: m.GetOrZeroValue(k)}
}

// GetValueByIndex get the value by index of key order.
//
// You should make sure 0 <= i < Len(), panic if out of bound.
func (m *Map[K, V]) GetValueByIndex(index int) V {
	k := m.GetKeyByIndex(index)
	return m.GetOrZeroValue(k)
}

func (m *Map[K, V]) set(key K, value V, alreadyExist bool) {
	if m.inner == nil {
		m.inner = make(map[K]V)
	}

	if !alreadyExist {
		m.order = append(m.order, key)
	}

	m.inner[key] = value
}

// Set a value by key without change its order, or place it at end if key is
// not exist.
//
// This operation is the same as [Map.Add] when duplicate key strategy is
// [UpdateValueKeepOrder].
func (m *Map[K, V]) Set(key K, value V) {
	m.set(key, value, m.Has(key))
}

// Add a key value pair.
//
// If the key is already exist in map, the behavior is controlled by
// [Map.DuplicateKeyStrategy].
func (m *Map[K, V]) Add(key K, value V) {
	var alreadyExist bool

	switch m.onDuplicateKey {
	default:
	case UpdateValueKeepOrder:
		{
			alreadyExist = m.Has(key)
		}
	case UpdateValueUpdateOrder:
		{
			m.Delete(key)
			// alreadyExist = false
		}
	case KeepValueUpdateOrder:
		{
			oldValue, exist := m.Get(key)
			if exist {
				value = oldValue
				m.Delete(key)
			}
			// alreadyExist = false
		}
	case Ignore:
		{
			if m.Has(key) {
				return
			}
			// alreadyExist = false
		}
	}

	m.set(key, value, alreadyExist)
}

// Append a series of kv pairs into map.
//
// The effect is consistent with calling [Map.Add](k, v) multi times.
func (m *Map[K, V]) Append(pairs ...Pair[K, V]) {
	for _, pair := range pairs {
		m.Add(pair.Key, pair.Value)
	}
}

// Delete a item by key.
//
// Performance: causes O(n) operation, avoid heavy use.
func (m *Map[K, V]) Delete(key K) {
	_, exist := m.inner[key]
	if !exist {
		return
	}

	for i, k := range m.order {
		if k == key {
			m.DeleteByIndex(i)
			return
		}
	}
}

// Delete a item by it's index in order.
//
// You should make sure 0 <= i < Len(), panic if out of bound.
//
// Performance: causes O(n) operation, avoid heavy use.
func (m *Map[K, V]) DeleteByIndex(index int) {
	key := m.order[index]
	m.order = append(m.order[:index], m.order[index+1:]...)
	delete(m.inner, key)
}

// Clean this map.
func (m *Map[K, V]) Clear() {
	m.order = nil
	m.inner = nil
}

// Len returns the size of map.
func (m *Map[K, V]) Len() int {
	return len(m.inner)
}

// Keys returns a copy of all keys of the map, in current order.
//
// Performance: O(n) operation. If you want iterate over the map,
// maybe [Map.Len] + [Map.GetKeyByIndex] is a better choice.
func (m *Map[K, V]) Keys() []K {
	// copy to avoid user modify the order.
	keys := make([]K, m.Len(), m.Len())
	copy(keys, m.order)
	return keys
}

// Values returns a copy of all values of the map, in current order.
//
// Performance: O(n) operation. If you want iterate over the map,
// maybe [Map.Len] + [Map.GetValueByIndex] is a better choice.
func (m *Map[K, V]) Values() []V {
	values := make([]V, 0, m.Len())
	for i := 0; i < m.Len(); i++ {
		values = append(values, m.GetValueByIndex(i))
	}
	return values
}

// Pairs gives you all data the map stored as a list of pair, in current order.
//
// Performance: O(n) operation. If you want iterate over the map,
// maybe [Map.Len] + [Map.GetByIndex] is a better choice.
func (m *Map[K, V]) Pairs() *PairList[K, V] {
	pairs := NewPairListWithCapacity[K, V](m.Len())

	for i := 0; i < m.Len(); i++ {
		pairs.List = append(pairs.List, m.GetByIndex(i))
	}

	return pairs
}

// Sort will reorder the map using the given less function.
func (m *Map[K, V]) Sort(lessFunc PairLessFunc[K, V]) {
	pairs := m.Pairs()

	pairs.Sort(lessFunc)

	for i := 0; i < m.Len(); i++ {
		m.order[i] = pairs.List[i].Key
	}
}

// Filter remove all item which make pred func return false.
//
// Performance: O(n) operation. More efficient then [Map.GetByIndex] +
// [Map.DeleteByIndex] in a loop, which is O(n^2).
func (m *Map[K, V]) Filter(pred PairFilterFunc[K, V]) {
	n := 0
	for i, length := 0, m.Len(); i < length; i++ {
		pair := m.GetByIndex(i)
		if pred(&pair) {
			m.order[n] = m.order[i]
			n++
		} else {
			delete(m.inner, pair.Key)
		}
	}
	m.order = m.order[:n]
}

// MarshalJSON implements [json.Marshaler] interface.
//
// You should not call this directly, use [json.Marshal] instead.
func (m Map[K, V]) MarshalJSON() ([]byte, error) {
	return marshalObject[K, V](&m)
}

// UnmarshalJSON implements [json.Unmarshaler] interface.
// You shouldn't call this directly, use [json.Unmarshal]/[JSONUnmarshal]
// instead.
func (m *Map[K, V]) UnmarshalJSON(data []byte) error {
	return unmarshalObject[K, V](data, m, OnDuplicateKey(m.onDuplicateKey))
}
