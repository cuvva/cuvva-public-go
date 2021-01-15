package snek

import (
	"unicode"

	stringcase "github.com/reiver/go-stringcase"
)

func DesnekKeysRecursive(in map[string]interface{}) map[string]interface{} {
	out := make(map[string]interface{}, len(in))

	for k, v := range in {
		if m, ok := v.(map[string]interface{}); ok {
			out[Desnek(k)] = DesnekKeysRecursive(m)
		} else {
			out[Desnek(k)] = v
		}
	}

	return out
}

func DesnekKeys(in map[string]interface{}) map[string]interface{} {
	out := make(map[string]interface{}, len(in))

	for k, v := range in {
		out[Desnek(k)] = v
	}

	return out
}

func Desnek(str string) string {
	return stringcase.ToCamelCase(str)
}

func SnekKeysRecursive(in map[string]interface{}) map[string]interface{} {
	out := make(map[string]interface{}, len(in))

	for k, v := range in {
		if m, ok := v.(map[string]interface{}); ok {
			out[Snek(k)] = SnekKeysRecursive(m)
		} else {
			out[Snek(k)] = v
		}
	}

	return out
}

func SnekKeys(in map[string]interface{}) map[string]interface{} {
	out := make(map[string]interface{}, len(in))

	for k, v := range in {
		out[Snek(k)] = v
	}

	return out
}

func Snek(str string) string {
	runes := []rune(str)
	length := len(runes)

	var out []rune
	for i := 0; i < length; i++ {
		if i > 0 && unicode.IsUpper(runes[i]) && ((i+1 < length && unicode.IsLower(runes[i+1])) || unicode.IsLower(runes[i-1])) {
			out = append(out, '_')
		}
		out = append(out, unicode.ToLower(runes[i]))
	}

	return string(out)
}
