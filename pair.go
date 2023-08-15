package geko

// Pair is k v pair.
type Pair[K, V any] struct {
	Key   K `json:"key"`
	Value V `json:"value"`
}

func CreatePair[K, V any](key K, value V) Pair[K, V] {
	return Pair[K, V]{key, value}
}

// PairLessFunc is the less func to sort a pair list.
type PairLessFunc[K, V any] func(a, b *Pair[K, V]) bool
