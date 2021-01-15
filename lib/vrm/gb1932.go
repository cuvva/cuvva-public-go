package vrm

import (
	"fmt"
)

type GB1932 struct {
	Reversed bool

	Serial   string
	Area     string
	Sequence string
}

func ParseGB1932(vrm string) VRM {
	if len(vrm) < 4 || len(vrm) > 6 {
		return nil
	}

	g := &GB1932{}

	if match(vrm[:3], isAlpha) {
		g.Serial = vrm[:1]
		g.Area = vrm[1:3]
		g.Sequence = vrm[3:]
	} else if match(vrm[:1], isNumeric) {
		g.Reversed = true

		area := vrm[len(vrm)-3:]
		if !match(area, isAlpha) {
			return nil
		}

		g.Serial = area[:1]
		g.Area = area[1:]

		g.Sequence = vrm[:len(vrm)-3]
	} else {
		return nil
	}

	if !match(g.Sequence, isNumeric) || any(g.Area, g.isProhibitedLetter) {
		return nil
	}

	return g
}

func (g *GB1932) Format() string {
	return "gb_1932"
}

func (g *GB1932) String() string {
	if g.Reversed {
		return fmt.Sprintf("%s%s%s", g.Sequence, g.Serial, g.Area)
	}

	return fmt.Sprintf("%s%s%s", g.Serial, g.Area, g.Sequence)
}

func (g *GB1932) PrettyString() string {
	if g.Reversed {
		return fmt.Sprintf("%s %s%s", g.Sequence, g.Serial, g.Area)
	}

	return fmt.Sprintf("%s%s %s", g.Serial, g.Area, g.Sequence)
}

func (g *GB1932) isProhibitedLetter(r rune) bool {
	return r == 'I' || r == 'Q' || r == 'Z'
}
