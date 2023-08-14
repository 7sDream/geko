package geko

// Wrapper type for a normal slice.
//
// Unmarshal from json into a List[any] will use our `Map[string]any` to store json object,
// instead of normal `map[string]any` from stdlib.
type List[T any] struct {
	List []T

	escapeHTML bool
}

func NewList[T any]() *List[T] {
	return NewListFrom[T](nil)
}

func NewListFrom[T any](list []T) *List[T] {
	return &List[T]{
		List:       list,
		escapeHTML: true,
	}
}

func NewListWithCapacity[T any](capacity int) *List[T] {
	return NewListFrom[T](make([]T, 0, capacity))
}

func (ko *List[T]) Get(i int) T {
	return ko.List[i]
}

func (ko *List[T]) Append(value ...T) {
	ko.List = append(ko.List, value...)
}

func (ko *List[T]) Len() int {
	return len(ko.List)
}

func (ko *List[T]) innerSlice() *[]T {
	return &ko.List
}

// MarshalJSON implements json.Marshaler interface.
// You should not call this directly, use json.Marshal(ko) instead.
func (kol List[T]) MarshalJSON() ([]byte, error) {
	return marshalArray[T](&kol)
}

// MarshalJSON implements json.Marshaler interface.
// You should not call this directly, use json.Marshal(ko) instead.
func (kol *List[T]) UnmarshalJSON(data []byte) error {
	return unmarshalArray[T](data, kol)
}
