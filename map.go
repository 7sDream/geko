package geko

// DuplicatedKeyStrategy controls the behavior of [Map.Add] when meet a
// duplicate key. Default strategy is [UpdateValueKeepOrder].
//
// If you want store all values of duplicated key, use [Pairs] type instead.
type DuplicatedKeyStrategy uint8

const (
	// UpdateValueKeepOrder will use new value, but keep the key order.
	//
	// {"a": 1, "b": 2, "a": 3} => {"a": 3, "b": 2}
	//
	// This is the default strategy.
	UpdateValueKeepOrder DuplicatedKeyStrategy = iota
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
// You can [Map.SetDuplicatedKeyStrategy] before call [json.Unmarshal] to
// control the behavior when object has duplicated key in your JSON string data.
//
// If you can't make sure the outmost item is object, see [JSONUnmarshal]
// function.
type Map[K comparable, V any] struct {
	order []K
	inner map[K]V

	duplicatedKeyStrategy DuplicatedKeyStrategy
}

// Object is [Map] whose type parameters are specialized as
// [string, any], used to represent dynamic objects in JSON.
type Object = *Map[string, any]

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

// DuplicatedKeyStrategy get current strategy when [Map.Add] with a duplicated
// key.
//
// See document of [DuplicatedKeyStrategy] and its enum value for detail.
func (m *Map[K, V]) DuplicatedKeyStrategy() DuplicatedKeyStrategy {
	return m.duplicatedKeyStrategy
}

// SetDuplicatedKeyStrategy set strategy when [Map.Add] with a duplicated key.
//
// See document of [DuplicatedKeyStrategy] and its enum value for detail.
func (m *Map[K, V]) SetDuplicatedKeyStrategy(strategy DuplicatedKeyStrategy) {
	m.duplicatedKeyStrategy = strategy
}

// Get a value by key. The second return value tells if the key exists. If
// not, returned value will be zero value of type V.
func (m *Map[K, V]) Get(key K) (V, bool) {
	v, exist := m.inner[key]
	return v, exist
}

// Has checks if key exist in the map.
func (m *Map[K, V]) Has(key K) bool {
	_, exist := m.inner[key]
	return exist
}

// GetOrZeroValue return stored value by key, or the zero value of type V
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
	return CreatePair(k, m.GetOrZeroValue(k))
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
// [Map.DuplicatedKeyStrategy].
func (m *Map[K, V]) Add(key K, value V) {
	var alreadyExist bool

	switch m.duplicatedKeyStrategy {
	default:
		fallthrough
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
			m.order = append(m.order[:i], m.order[i+1:]...)
			break
		}
	}

	delete(m.inner, key)
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
	length := m.Len()

	values := make([]V, 0, length)
	for i := 0; i < length; i++ {
		values = append(values, m.GetValueByIndex(i))
	}

	return values
}

// Pairs gives you all data the map stored as a list of pair, in current order.
//
// Performance: O(n) operation. If you want iterate over the map,
// maybe [Map.Len] + [Map.GetByIndex] is a better choice.
func (m *Map[K, V]) Pairs() *Pairs[K, V] {
	length := m.Len()

	pairs := NewPairsWithCapacity[K, V](length)

	for i := 0; i < length; i++ {
		pairs.List = append(pairs.List, m.GetByIndex(i))
	}

	return pairs
}

// Sort will reorder the map using the given less function.
func (m *Map[K, V]) Sort(lessFunc PairLessFunc[K, V]) {
	pairs := m.Pairs()

	pairs.Sort(lessFunc)

	for i, length := 0, m.Len(); i < length; i++ {
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
//
// You shouldn't call this directly, use [json.Unmarshal]/[JSONUnmarshal]
// instead.
func (m *Map[K, V]) UnmarshalJSON(data []byte) error {
	return unmarshalObject[K, V](
		data, m,
		UseObject(),
		ObjectOnDuplicatedKey(m.duplicatedKeyStrategy),
	)
}
