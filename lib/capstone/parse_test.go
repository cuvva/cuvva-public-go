package capstone_test

import (
	"errors"
	"testing"

	"github.com/cuvva/cuvva-public-go/lib/capstone"
	"github.com/cuvva/cuvva-public-go/lib/capstone/decoder"
)

// Tests if the struct tag exists on the provided type
type MissingStructTag struct {
	Test *string
}

func (mst MissingStructTag) Format() capstone.Format {
	return capstone.Format{
		{
			Offset:   0,
			Length:   1,
			TagValue: "where_am_i",
			Decoder:  decoder.Int{},
		},
	}
}

func TestMissingStructTag(t *testing.T) {
	var mst MissingStructTag

	err := capstone.Parse("1", &mst)
	if !errors.Is(err, capstone.ErrCapstoneTagMissing) {
		t.Errorf("expected error '%s', got '%s'", capstone.ErrCapstoneTagMissing, err)
	}
}

type ExampleCapstone struct {
	SomeValue string `cap:"some_value"`
}

func (ec ExampleCapstone) Format() capstone.Format {
	return capstone.Format{
		{
			Offset:   0,
			Length:   1,
			TagValue: "some_value",
			Decoder: decoder.String{
				UnavailableValue: "Z",
				Values: map[string]string{
					"F": "respect",
				},
			},
		},
	}
}

// Test if the string provided is too short to populate the given capstone struct
func TestCapstoneTooShort(t *testing.T) {
	var ec ExampleCapstone

	err := capstone.Parse("", &ec)
	if !errors.Is(err, capstone.ErrInputTooShort) {
		t.Errorf("expected error '%s', got '%s'", capstone.ErrInputTooShort, err)
	}
}

// Tests if the provided capstone struct is not a pointer
func TestCapstoneNotPtr(t *testing.T) {
	var ec ExampleCapstone

	err := capstone.Parse("", ec)
	if !errors.Is(err, capstone.ErrCapstoneNotPtr) {
		t.Errorf("expected error '%s', got '%s'", capstone.ErrCapstoneNotPtr, err)
	}
}

type BadTagExample ExampleCapstone

func (bte BadTagExample) Format() capstone.Format {
	return capstone.Format{
		{
			Offset:   0,
			Length:   1,
			TagValue: "x",
			Decoder: decoder.Int{
				UnavailableValue: "X",
			},
		},
	}
}

// Tests if the parser errors correctly if a format is provided where there's
// no appropriate tag to apply in the struct
func TestBadTag(t *testing.T) {
	var bte BadTagExample

	err := capstone.Parse("1", &bte)
	if !errors.Is(err, capstone.ErrInvalidTag) {
		t.Errorf("expected error '%s', got '%s'", capstone.ErrInvalidTag, err)
	}
}

type UnsupportedTypeExample struct {
	UnsupportedField *struct{} `cap:"some_value"`
}

type UnsupportedDecoder struct{}

func (ud UnsupportedDecoder) Decode(in string) (interface{}, error) {
	return &struct{}{}, nil
}

func (bte UnsupportedTypeExample) Format() capstone.Format {
	return capstone.Format{
		{
			Offset:   0,
			Length:   1,
			TagValue: "some_value",
			Decoder:  UnsupportedDecoder{},
		},
	}
}

// Tests if the decoded value type is settable by the parser
func TestUnsupportedDecoder(t *testing.T) {
	var ute UnsupportedTypeExample

	err := capstone.Parse("x", &ute)
	if !errors.Is(err, capstone.ErrUnsupportedTypeSet) {
		t.Errorf("expected error '%s', got '%s'", capstone.ErrUnsupportedTypeSet, err)
	}
}

type NonPointerExample struct {
	UnsupportedField struct{} `cap:"some_value"`
}

type NonPointerDecoder struct{}

func (ud NonPointerDecoder) Decode(in string) (interface{}, error) {
	return struct{}{}, nil
}

func (npe NonPointerExample) Format() capstone.Format {
	return capstone.Format{
		{
			Offset:   0,
			Length:   1,
			TagValue: "some_value",
			Decoder:  NonPointerDecoder{},
		},
	}
}

// Tests to ensure the value deocded is infact a pointer
func TestValueNotPointer(t *testing.T) {
	var npd NonPointerExample

	err := capstone.Parse("x", &npd)
	if !errors.Is(err, capstone.ErrExpectedDecodedPtr) {
		t.Errorf("expected error '%s', got '%s'", capstone.ErrExpectedDecodedPtr, err)
	}
}

type MismatchTypeExample ExampleCapstone
type MismatchTypeDecoder struct{}

func (mtd MismatchTypeDecoder) Decode(in string) (interface{}, error) {
	return 1, nil
}

func (mte MismatchTypeExample) Format() capstone.Format {
	return capstone.Format{
		{
			Offset:   0,
			Length:   1,
			TagValue: "some_value",
			Decoder:  MismatchTypeDecoder{},
		},
	}
}

// Ensures that the parser errors if the type returned by the decoder does
// not match that of the struct field.
func TestMismatchType(t *testing.T) {
	var tmt MismatchTypeExample

	err := capstone.Parse("-", &tmt)
	if !errors.Is(err, capstone.ErrDecodedTypeMismatch) {
		t.Errorf("expected error '%s', got '%s'", capstone.ErrDecodedTypeMismatch, err)
	}
}
