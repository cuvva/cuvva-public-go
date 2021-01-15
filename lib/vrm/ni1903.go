package vrm

import (
	"fmt"
)

type NI1903 struct {
	Reversed bool

	Area     string
	Sequence string
}

func ParseNI1903(vrm string) VRM {
	if len(vrm) < 3 || len(vrm) > 6 {
		return nil
	}

	n := &NI1903{}

	if match(vrm[:2], isAlpha) {
		n.Area = vrm[:2]
		n.Sequence = vrm[2:]
	} else if match(vrm[:1], isNumeric) {
		n.Reversed = true
		n.Area = vrm[len(vrm)-2:]
		if !match(n.Area, isAlpha) {
			return nil
		}

		n.Sequence = vrm[:len(vrm)-2]
	} else {
		return nil
	}

	if !match(n.Sequence, isNumeric) || !n.acceptableArea(n.Area) {
		return nil
	}

	return n
}

func (n *NI1903) Format() string {
	return "ni_1903"
}

func (n *NI1903) String() string {
	if n.Reversed {
		return fmt.Sprintf("%s%s", n.Sequence, n.Area)
	}

	return fmt.Sprintf("%s%s", n.Area, n.Sequence)
}

func (n *NI1903) PrettyString() string {
	if n.Reversed {
		return fmt.Sprintf("%s %s", n.Sequence, n.Area)
	}

	return fmt.Sprintf("%s %s", n.Area, n.Sequence)
}

func (n *NI1903) acceptableArea(area string) bool {
	for _, r := range area {
		if r == 'I' || r == 'Z' {
			return true
		}
	}

	return false
}
