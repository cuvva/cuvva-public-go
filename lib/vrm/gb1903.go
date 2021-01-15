package vrm

import (
	"fmt"
)

type GB1903 struct {
	Reversed bool

	Area     string
	Sequence string
}

func ParseGB1903(vrm string) VRM {
	if len(vrm) < 2 || len(vrm) > 6 {
		return nil
	}

	g := &GB1903{}

	if match(vrm[:1], isAlpha) {
		if match(vrm[:2], isAlpha) {
			g.Area = vrm[:2]
			g.Sequence = vrm[2:]
		} else {
			g.Area = vrm[:1]
			g.Sequence = vrm[1:]
		}
	} else if match(vrm[:1], isNumeric) {
		g.Reversed = true

		if match(vrm[len(vrm)-2:], isAlpha) {
			g.Area = vrm[len(vrm)-2:]
			g.Sequence = vrm[:len(vrm)-2]
		} else {
			g.Area = vrm[len(vrm)-1:]
			g.Sequence = vrm[:len(vrm)-1]
		}
	} else {
		return nil
	}

	if !match(g.Sequence, isNumeric) || any(g.Area, g.isProhibitedLetter) {
		return nil
	}

	return g
}

func (g *GB1903) Format() string {
	return "gb_1903"
}

func (g *GB1903) String() string {
	if g.Reversed {
		return fmt.Sprintf("%s%s", g.Sequence, g.Area)
	}

	return fmt.Sprintf("%s%s", g.Area, g.Sequence)
}

func (g *GB1903) PrettyString() string {
	if g.Reversed {
		return fmt.Sprintf("%s %s", g.Sequence, g.Area)
	}

	return fmt.Sprintf("%s %s", g.Area, g.Sequence)
}

func (g *GB1903) isProhibitedLetter(r rune) bool {
	return r == 'I' || r == 'Q' || r == 'Z'
}
