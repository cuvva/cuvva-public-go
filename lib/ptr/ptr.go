package ptr

import (
	"time"
)

func Ptr[T any](v T) *T {
	return &v
}

func String(v string) *string {
	return &v
}

func Bool(v bool) *bool {
	return &v
}

func Int(v int) *int {
	return &v
}

func Int64(v int64) *int64 {
	return &v
}

func Float64(v float64) *float64 {
	return &v
}

func Time(v time.Time) *time.Time {
	return &v
}

// Equal returns true if the two parameters are both nil or are pointers to the same value. Provides a type-safe generic
// alternative to reflect.DeepEqual that ensures both parameters are pointers to the same comparable type.
func Equal[T comparable](a *T, b *T) bool {
	if a == nil || b == nil {
		return a == b
	}

	return *a == *b
}

// Copy returns nil if the given pointer is nil, otherwise it copies the pointed-to value into a new variable and
// returns a pointer to that copy. This is useful when assigning pointer fields between structs so that mutating one
// struct does not affect the other.
func Copy[T any](v *T) *T {
	if v == nil {
		return nil
	}

	c := *v
	return &c
}
