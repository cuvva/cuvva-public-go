package capstone

import (
	"errors"
)

// Decoder interface used for all element decoders
type Decoder interface {
	Decode(in string) (interface{}, error)
}

// Format is a slice of elements
type Format []Element

// Element holds the position info of each capstone attribute, the name which
// should be used to populate the struct tag, and a decoder which implements
// Decoder and handles decoding this element type.
type Element struct {
	Offset   int
	Length   int
	TagValue string
	Decoder  Decoder
}

// Capstone interface ensures all provided types to Parse have a valid Format
type Capstone interface {
	Format() Format
}

// ErrCapstoneNotPtr is returned if the capstone provided is not a pointer is nil
var ErrCapstoneNotPtr = errors.New("capstone given is not a pointer or is nil")

// ErrInputTooShort is used if the string provided is too short to process a
// whole capstone
var ErrInputTooShort = errors.New("provided value too short for capstone format")

// ErrCapstoneTagMissing is returned if a `cap` struct tag is missing on a capstone
// struct field
var ErrCapstoneTagMissing = errors.New("missing 'cap' struct tag on capstone")

// ErrDecodedTypeMismatch happens when a decoded value cannot be set to the
// destination struct as the underlying types do not match
var ErrDecodedTypeMismatch = errors.New("decoded value does not match type in struct")

// ErrExpectedDecodedPtr is thrown when a decoded value to be set on a struct is not
// a pointer. All decoded values should be pointers as all values are optional.
var ErrExpectedDecodedPtr = errors.New("expected a pointer value for decoding field")

// ErrUnsupportedTypeSet is returned when the decoded value being set on the struct
// provided is not 'settable', i.e. is not yet supported to be set.
var ErrUnsupportedTypeSet = errors.New("unsupported type to unmarshal")

// ErrInvalidTag is used to indicate that the given Format element tag does not
// match any of the provided struct
var ErrInvalidTag = errors.New("invalid struct tag")
