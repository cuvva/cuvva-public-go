package vrm

import (
	"fmt"
	"strings"
	"unicode"
)

// VRM represents a single coerced VRM.
type VRM interface {
	// Format returns the name of the VRM scheme identified.
	Format() string

	// String returns a stringified version of the normalised VRM.
	String() string

	// PrettyString returns a stringified version of the VRM as seen.
	PrettyString() string
}

// Parser takes a normalised VRM and returns a structured implementing the VRM interface,
// or nil if the given VRM could not be parsed by this format.
type Parser func(vrm string) VRM

// DefaultParsers is the default set of all parsers implemented by this package, ordered
// by likelyhood of appearance.
var DefaultParsers = []Parser{
	// current schemes
	ParseGB2001,
	ParseNI1966,
	ParseMilitary,
	ParseDiplomatic,

	// historic schemes
	ParseGB1983,
	ParseGB1963,
	ParseGB1932,
	ParseGB1903,
	ParseNI1903,
}

// Coerces the input into a set of possible VRMs which the input could represent. The returned
// array contains the VRM details for each given format, sorted in order of likelihood, where
// the most likely format is the first value. If the allowed formats are specified, coercion
// will be limited to these formats. Any other formats will not be checked.
//
// returns an empty array if the input is invalid or cannot be coerced into any of
// the formats checked.
func Coerce(vrm string, formats ...Parser) (results []VRM) {
	if len(formats) == 0 {
		formats = DefaultParsers
	}

	vrm = NormaliseVRM(vrm)
	if !ValidVRM(vrm) {
		return
	}

	// apply parsers to given normalised vrm
	results = applyParser(vrm, formats...)

	// if no results found, apply parsers to variations on normalised vrm
	variants := combinations("", vrm, substitutions)

	for _, variant := range variants {
		results = append(results, applyParser(variant, formats...)...)
	}

	return
}

func applyParser(vrm string, formats ...Parser) (results []VRM) {
	for _, format := range formats {
		result := format(vrm)
		if result != nil {
			checkIntegrity(vrm, result)
			results = append(results, result)
		}
	}

	return
}

func checkIntegrity(input string, result VRM) {
	strippedPretty := strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, result.PrettyString())

	if result.String() != input || strippedPretty != input {
		panic(fmt.Sprintf("input: '%s', output: '%s', pretty: '%s'", input, result.String(), result.PrettyString()))
	}
}

// Verifies that the given VRM matches one of the known formats and returns the relevant VRM
// details. Only normalized VRMs are accepted. If the format is specified, only that format is checked.
//
// returns nil if the VRM does not match any of the formats checked.
func Info(vrm string, formats ...Parser) VRM {
	if len(formats) == 0 {
		formats = DefaultParsers
	}

	if !ValidVRM(vrm) {
		return nil
	}

	vrms := applyParser(vrm, formats...)
	if len(vrms) > 0 {
		return vrms[0]
	}

	return nil
}

// return true if given vrm is valid (for the formats implements in this package).
func ValidVRM(vrm string) bool {
	if len(vrm) < 2 || len(vrm) > 7 {
		return false
	}

	for _, r := range vrm {
		if !isAlpha(r) && !isNumeric(r) {
			return false
		}
	}

	return true
}

// return the given VRM normalised (to [A-Z0-9])
func NormaliseVRM(vrm string) string {
	nvrm := make([]byte, len(vrm))
	nlen := 0

	for _, r := range vrm {
		if isAlpha(r) {
			nvrm[nlen] = byte(r)
			nlen++
		} else if r >= 'a' && r <= 'z' {
			nvrm[nlen] = byte(r - 32)
			nlen++
		} else if isNumeric(r) {
			nvrm[nlen] = byte(r)
			nlen++
		}
	}

	return string(nvrm[:nlen])
}

// return true if string contains only runes defined by function
func match(s string, fn func(rune) bool) bool {
	if len(s) == 0 {
		return false
	}

	for _, r := range s {
		if !fn(r) {
			return false
		}
	}

	return true
}

func any(s string, fn func(rune) bool) bool {
	if len(s) == 0 {
		return false
	}

	for _, r := range s {
		if fn(r) {
			return true
		}
	}

	return false
}

// return true if rune if A-Z
func isAlpha(r rune) bool {
	return r >= 'A' && r <= 'Z'
}

// return true if rune is 0-9
func isNumeric(r rune) bool {
	return r >= '0' && r <= '9'
}

// substitutions applies the following substitutions:
// 1 > I
// I > 1
// 0 > O
// O > 0
func substitutions(r rune) rune {
	if r == '1' {
		return 'I'
	} else if r == 'I' {
		return '1'
	} else if r == '0' {
		return 'O'
	} else if r == 'O' {
		return '0'
	}

	return r
}

// combinations applies a substitution function to a string (with optional prefix)
// and returns all possible variations.
func combinations(prefix, s string, sub func(rune) rune) (res []string) {
	for i, r := range s {
		if b := sub(r); b != r {
			res = append(res, fmt.Sprintf("%s%s%c%s", prefix, s[:i], b, s[i+1:]))

			variants := combinations(fmt.Sprintf("%s%s%c", prefix, s[:i], b), s[i+1:], sub)

			if len(variants) > 0 {
				res = append(res, variants...)
			}
		}
	}

	return
}
