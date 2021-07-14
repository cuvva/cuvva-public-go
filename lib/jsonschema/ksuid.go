package jsonschema

import (
	"github.com/cuvva/cuvva-public-go/lib/ksuid"
	"github.com/xeipuuv/gojsonschema"
)

func init() {
	gojsonschema.FormatCheckers.Add("ksuid", ksuidFormatChecker{})
}

type ksuidFormatChecker struct{}

func (f ksuidFormatChecker) IsFormat(input interface{}) bool {
	str, ok := input.(string)
	if !ok {
		return false
	}

	_, err := ksuid.Parse(str)
	return err == nil
}
