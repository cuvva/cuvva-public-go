package vrm

import (
	"fmt"
)

type GB1963 struct {
	Serial   string
	Area     string
	Sequence string
	AgeID    string
}

func ParseGB1963(vrm string) VRM {
	if len(vrm) < 5 || len(vrm) > 7 {
		return nil
	}

	area := vrm[:3]
	if !match(area, isAlpha) {
		return nil
	}

	g := &GB1963{
		Serial: area[:1],
		Area:   area[1:],
	}

	if any(area, g.isProhibitedLetter) {
		return nil
	}

	g.Sequence = vrm[3 : len(vrm)-1]
	if !match(g.Sequence, isNumeric) {
		return nil
	}

	g.AgeID = vrm[len(vrm)-1:]
	if _, ok := ageIDToYear63[g.AgeID[0]]; !ok {
		return nil
	}

	return g
}

func (g *GB1963) Format() string {
	return "gb_1963"
}

func (g *GB1963) String() string {
	return fmt.Sprintf("%s%s%s%s", g.Serial, g.Area, g.Sequence, g.AgeID)
}

func (g *GB1963) PrettyString() string {
	return fmt.Sprintf("%s%s %s%s", g.Serial, g.Area, g.Sequence, g.AgeID)
}

func (g *GB1963) isProhibitedLetter(r rune) bool {
	return r == 'I' || r == 'Q' || r == 'Z'
}

var ageIDToYear63 = map[byte]int{
	'A': 1963,
	'B': 1964,
	'C': 1965,
	'D': 1966,
	'E': 1967,
	'F': 1968,
	'G': 1969,
	'H': 1970,
	'J': 1971,
	'K': 1972,
	'L': 1973,
	'M': 1974,
	'N': 1975,
	'P': 1976,
	'R': 1977,
	'S': 1978,
	'T': 1979,
	'V': 1980,
	'W': 1981,
	'X': 1982,
	'Y': 1983,
}
