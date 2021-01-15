package capstone

import (
	"fmt"
	"reflect"
)

// Parse takes the raw capstone string, parses and unmarshals it into cap.
// - cap must be a pointer
func Parse(v string, cap Capstone) error {
	rv := reflect.ValueOf(cap)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return ErrCapstoneNotPtr
	}

	for _, el := range cap.Format() {
		if el.Offset+el.Length > len(v) {
			return fmt.Errorf("%w: input was %d chars in length", ErrInputTooShort, len(v))
		}

		p := v[el.Offset : el.Length+el.Offset]

		v, err := el.Decoder.Decode(p)
		if err != nil {
			return err
		}

		if err := populateField(v, el, cap); err != nil {
			return err
		}
	}

	return nil
}

func populateField(v interface{}, el Element, cap Capstone) error {
	if v == nil {
		return nil
	}

	fieldVal, fieldName, err := getValueForTag(cap, "cap", el.TagValue)
	if err != nil {
		return err
	}

	if !fieldVal.CanSet() {
		return fmt.Errorf("cannot set capstone value at field %s", fieldName)
	}

	valType, fieldType := getType(v), getType(fieldVal.Interface())
	if valType != fieldType {
		return fmt.Errorf("%w: type %s used for field %s (%s)", ErrDecodedTypeMismatch, valType, fieldType, fieldName)
	}

	val := reflect.ValueOf(v)

	// NOTE(sn): all capstone fields are optional and as such we expect that
	// all target structs contain only pointer fields.
	if val.Kind() != reflect.Ptr {
		return fmt.Errorf("%w: %s", ErrExpectedDecodedPtr, fieldName)
	}

	switch valType {
	case reflect.String.String():
		fieldVal.Set(val)
	case reflect.Int.String():
		fieldVal.Set(val)
	case reflect.Float64.String():
		fieldVal.Set(val)
	case reflect.Bool.String():
		fieldVal.Set(val)
	default:
		return fmt.Errorf("%w: attempted to set %s for field %s", ErrUnsupportedTypeSet, valType, fieldName)
	}

	return nil
}

// getValueForTag returns a reflect.Value and field name, matching a given
// struct tag and value found on any field from v. Will return an error if not
// found, or an error if the tag does not appear on the struct.
func getValueForTag(v interface{}, tag, lookup string) (reflect.Value, string, error) {
	val := reflect.ValueOf(v).Elem()
	t := val.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		key, ok := field.Tag.Lookup(tag)
		if !ok {
			return reflect.Value{}, "", fmt.Errorf("%w: missing tag on field %s", ErrCapstoneTagMissing, field.Name)
		}

		if key != lookup {
			continue
		}

		return val.Field(i), field.Name, nil
	}

	return reflect.Value{}, "", fmt.Errorf("%w: no value %s found for tag %s", ErrInvalidTag, lookup, tag)
}

// getType returns a string of the concrete underlying type, for example if
// *string, will return string.
func getType(v interface{}) string {
	t := reflect.TypeOf(v)

	if t.Kind() == reflect.Ptr {
		return t.Elem().Name()
	}

	return t.Name()
}
