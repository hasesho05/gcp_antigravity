package util

// ToPointer returns a pointer to the given value.
func ToPointer[T any](v T) *T {
	return &v
}

// FromPointer returns the value pointed to by v, or the zero value of T if v is nil.
func FromPointer[T any](v *T) T {
	if v == nil {
		var zero T
		return zero
	}
	return *v
}
