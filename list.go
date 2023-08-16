package geko

// Wrapper type for a normal slice.
//
// Unmarshal from JSON into a *List[any] will use *[Map][string, any] from this
// package to store JSON object, instead of normal map[string]any from std lib.
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
func (kol *List[T]) Get(index int) T {
	return kol.List[index]
}

// Set value at index.
func (kol *List[T]) Set(index int, value T) {
	kol.List[index] = value
}

// Append values into list.
func (kol *List[T]) Append(value ...T) {
	kol.List = append(kol.List, value...)
}

// Delete value at index.
func (kol *List[T]) Delete(index int) {
	kol.List = append(kol.List[:index], kol.List[index+1:]...)
}

// Len give length of the list.
func (kol *List[T]) Len() int {
	return len(kol.List)
}

func (kol *List[T]) innerSlice() *[]T {
	return &kol.List
}

// MarshalJSON implements [json.Marshaler] interface.
// You should not call this directly, use [json.Marshal] instead.
func (kol List[T]) MarshalJSON() ([]byte, error) {
	return marshalArray[T](&kol)
}

// UnmarshalJSON implements [json.Unmarshaler] interface.
// You should not call this directly, use [json.Marshal] instead.
func (kol *List[T]) UnmarshalJSON(data []byte) error {
	return unmarshalArray[T](data, kol)
}
