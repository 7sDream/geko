package geko

// Wrapper type for a normal slice.
//
// Unmarshal from JSON into a *[List][any] will use *[Map][string, any] from
// this package to store JSON object, use *[List][any] to store JSON array,
// instead of normal map[string]any and []any from std lib.
type List[T any] struct {
	List []T
}

// NewList create a new empty List.
func NewList[T any]() *List[T] {
	return NewListFrom[T](nil)
}

// NewListFrom create a List from a slice.
func NewListFrom[T any](list []T) *List[T] {
	return &List[T]{
		List: list,
	}
}

// NewList create a new empty List, but init with some capacity, for optimize
// memory allocation.
func NewListWithCapacity[T any](capacity int) *List[T] {
	return NewListFrom[T](make([]T, 0, capacity))
}

// Get value at index.
func (l *List[T]) Get(index int) T {
	return l.List[index]
}

// Set value at index.
func (l *List[T]) Set(index int, value T) {
	l.List[index] = value
}

// Append values into list.
func (l *List[T]) Append(value ...T) {
	l.List = append(l.List, value...)
}

// Delete value at index.
func (l *List[T]) Delete(index int) {
	l.List = append(l.List[:index], l.List[index+1:]...)
}

// Len give length of the list.
func (l *List[T]) Len() int {
	return len(l.List)
}

func (l *List[T]) innerSlice() *[]T {
	return &l.List
}

// MarshalJSON implements [json.Marshaler] interface.
// You should not call this directly, use [json.Marshal] instead.
func (l List[T]) MarshalJSON() ([]byte, error) {
	return marshalArray[T](&l)
}

// UnmarshalJSON implements [json.Unmarshaler] interface.
// You should not call this directly, use [json.Marshal] instead.
func (l *List[T]) UnmarshalJSON(data []byte) error {
	return unmarshalArray[T](data, l)
}
