package vrm

import (
	"fmt"
)

type GB1983 struct {
	AgeID    string
	Sequence string
	Serial   string
	Area     string
}

func ParseGB1983(vrm string) VRM {
	if len(vrm) < 5 || len(vrm) > 7 {
		return nil
	}

	area := vrm[len(vrm)-3:]
	if !match(area, isAlpha) {
		return nil
	}

	g := &GB1983{
		Serial: area[:1],
		Area:   area[1:],
	}

	if match(area, g.isProhibitedLetter) {
		return nil
	}

	g.AgeID = vrm[:1]
	if _, ok := ageIdToYear83[g.AgeID[0]]; !ok {
		return nil
	}

	g.Sequence = vrm[1 : len(vrm)-3]
	if !match(g.Sequence, isNumeric) {
		return nil
	}

	return g
}

func (g *GB1983) Format() string {
	return "gb_1983"
}

func (g *GB1983) String() string {
	return fmt.Sprintf("%s%s%s%s", g.AgeID, g.Sequence, g.Serial, g.Area)
}

func (g *GB1983) PrettyString() string {
	return fmt.Sprintf("%s%s %s%s", g.AgeID, g.Sequence, g.Serial, g.Area)
}

func (g *GB1983) isProhibitedLetter(r rune) bool {
	return r == 'I' || r == 'Q' || r == 'Z'
}

// DANGER: HERE BE DRAGONS
// The mapping does not contain sequential letters nor incremeting years
var ageIdToYear83 = map[byte]int{
	'A': 1984,
	'B': 1985,
	'C': 1986,
	'D': 1987,
	'E': 1988,
	'F': 1989,
	'G': 1990,
	'H': 1991,
	'J': 1992,
	'K': 1993,
	'L': 1994,
	'M': 1995,
	'N': 1996,
	'P': 1997,
	'R': 1998,
	'S': 1999,
	'T': 1999,
	'V': 2000,
	'W': 2000,
	'X': 2001,
	'Y': 2001,
}
