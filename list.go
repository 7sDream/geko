package geko

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

func (ko *List[T]) SetEscapeHTML(escape bool) {
	ko.escapeHTML = escape
}

func (ko *List[T]) EscapeHTML() bool {
	return ko.escapeHTML
}

func (ko *List[T]) innerSlice() *[]T {
	return &ko.List
}

func (ko *List[T]) append(value T) {
	ko.List = append(ko.List, value)
}

// MarshalJSON implements json.Marshaler interface.
// You should not call this directly, use json.Marshal(ko) instead.
func (kol List[T]) MarshalJSON() ([]byte, error) {
	return marshalArray(&kol)
}

// MarshalJSON implements json.Marshaler interface.
// You should not call this directly, use json.Marshal(ko) instead.
func (kol *List[T]) UnmarshalJSON(data []byte) error {
	return unmarshalArray(data, kol)
}
