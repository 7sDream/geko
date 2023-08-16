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
	kom := NewMap[K, V]()
	kom.order = make([]K, 0, capacity)
	kom.inner = make(map[K]V, capacity)
	return kom
}

// DuplicateKeyStrategy get current strategy when Set with a duplicate key.
//
// See document of [DuplicateKeyStrategy] and its enum value for detail.
func (kom *Map[K, V]) DuplicateKeyStrategy() DuplicateKeyStrategy {
	return kom.onDuplicateKey
}

// SetDuplicateKeyStrategy set strategy when [Map.Set] with a duplicate key.
//
// See document of [DuplicateKeyStrategy] and its enum value for detail.
func (kom *Map[K, V]) SetDuplicateKeyStrategy(strategy DuplicateKeyStrategy) {
	kom.onDuplicateKey = strategy
}

// Get a value by key. The second return value is true if the key exists,
// otherwise false.
func (kom *Map[K, V]) Get(key K) (V, bool) {
	v, exist := kom.inner[key]
	return v, exist
}

// Has checks if key is in the map.
func (kom *Map[K, V]) Has(key K) bool {
	_, exist := kom.inner[key]
	return exist
}

// GetOrZeroValue return stored value by key, or the zero value of value type
// if key not exist.
func (kom *Map[K, V]) GetOrZeroValue(key K) V {
	return kom.inner[key]
}

// GetKeyByIndex get key by it's index in order.
//
// You should make sure 0 <= i < Len(), panic if out of bound.
func (kom *Map[K, V]) GetKeyByIndex(index int) K {
	return kom.order[index]
}

// GetByIndex get the key and value by index of key order.
//
// You should make sure 0 <= i < Len(), panic if out of bound.
func (kom *Map[K, V]) GetByIndex(index int) Pair[K, V] {
	k := kom.GetKeyByIndex(index)
	return Pair[K, V]{Key: k, Value: kom.GetOrZeroValue(k)}
}

// GetValueByIndex get the value by index of key order.
//
// You should make sure 0 <= i < Len(), panic if out of bound.
func (kom *Map[K, V]) GetValueByIndex(index int) V {
	k := kom.GetKeyByIndex(index)
	return kom.GetOrZeroValue(k)
}

func (kom *Map[K, V]) set(key K, value V, alreadyExist bool) {
	if kom.inner == nil {
		kom.inner = make(map[K]V)
	}

	if !alreadyExist {
		kom.order = append(kom.order, key)
	}

	kom.inner[key] = value
}

// Set a value by key.
//
// If the key is already exist in map, the behavior is controlled by
// [Map.DuplicateKeyStrategy].
func (kom *Map[K, V]) Set(key K, value V) {
	var alreadyExist bool

	switch kom.onDuplicateKey {
	default:
	case UpdateValueKeepOrder:
		{
			alreadyExist = kom.Has(key)
		}
	case UpdateValueUpdateOrder:
		{
			kom.Delete(key)
			// alreadyExist = false
		}
	case KeepValueUpdateOrder:
		{
			oldValue, exist := kom.Get(key)
			if exist {
				value = oldValue
				kom.Delete(key)
			}
			// alreadyExist = false
		}
	case Ignore:
		{
			if kom.Has(key) {
				return
			}
			// alreadyExist = false
		}
	}

	kom.set(key, value, alreadyExist)
}

// Append a series of kv pairs into map.
//
// The effect is consistent with calling [Map.Set](k, v) multi times.
func (kom *Map[K, V]) Append(pairs ...Pair[K, V]) {
	for _, pair := range pairs {
		kom.Set(pair.Key, pair.Value)
	}
}

// Delete a item by key.
//
// Performance: causes O(n) operation, avoid heavy use.
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
//
// You should make sure 0 <= i < Len(), panic if out of bound.
//
// Performance: causes O(n) operation, avoid heavy use.
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

// Keys returns a copy of all keys of the map, in current order.
//
// Performance: O(n) operation. If you want iterate over the map,
// maybe [Map.Len] + [Map.GetKeyByIndex] is a better choice.
func (kom *Map[K, V]) Keys() []K {
	// copy to avoid user modify the order.
	keys := make([]K, kom.Len(), kom.Len())
	copy(keys, kom.order)
	return keys
}

// Values returns a copy of all values of the map, in current order.
//
// Performance: O(n) operation. If you want iterate over the map,
// maybe [Map.Len] + [Map.GetValueByIndex] is a better choice.
func (kom *Map[K, V]) Values() []V {
	values := make([]V, 0, kom.Len())
	for i := 0; i < kom.Len(); i++ {
		values = append(values, kom.GetValueByIndex(i))
	}
	return values
}

// Pairs gives you all data the map stored as a list of pair, in current order.
//
// Performance: O(n) operation. If you want iterate over the map,
// maybe [Map.Len] + [Map.GetByIndex] is a better choice.
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

// Filter remove all item which make pred func return false.
//
// Performance: O(n) operation. More efficient then [Map.GetByIndex] +
// [Map.DeleteByIndex] in a loop, which is O(n^2).
func (kom *Map[K, V]) Filter(pred PairFilterFunc[K, V]) {
	n := 0
	for i, length := 0, kom.Len(); i < length; i++ {
		pair := kom.GetByIndex(i)
		if pred(&pair) {
			kom.order[n] = kom.order[i]
			n++
		} else {
			delete(kom.inner, pair.Key)
		}
	}
	kom.order = kom.order[:n]
}

// MarshalJSON implements [json.Marshaler] interface.
//
// You should not call this directly, use [json.Marshal] instead.
func (kom Map[K, V]) MarshalJSON() ([]byte, error) {
	return marshalObject[K, V](&kom)
}

// UnmarshalJSON implements [json.Unmarshaler] interface.
// You shouldn't call this directly, use [json.Unmarshal]/[JSONUnmarshal]
// instead.
func (kom *Map[K, V]) UnmarshalJSON(data []byte) error {
	return unmarshalObject[K, V](data, kom, OnDuplicateKey(kom.onDuplicateKey))
}
