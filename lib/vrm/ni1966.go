package vrm

import (
	"fmt"
)

type NI1966 struct {
	Serial   string
	Area     string
	Sequence string
}

func ParseNI1966(vrm string) VRM {
	if len(vrm) < 4 || len(vrm) > 7 {
		return nil
	}

	n := &NI1966{
		Serial: vrm[:1],
		Area:   vrm[1:3],
	}

	if !match(n.Serial, isAlpha) {
		return nil
	} else if !match(n.Area, isAlpha) || !n.acceptableArea(n.Area) {
		return nil
	}

	n.Sequence = vrm[3:]
	if !match(n.Sequence, isNumeric) {
		return nil
	}

	return n
}

func (n *NI1966) Format() string {
	return "ni_1966"
}

func (n *NI1966) String() string {
	return fmt.Sprintf("%s%s%s", n.Serial, n.Area, n.Sequence)
}

func (n *NI1966) PrettyString() string {
	return fmt.Sprintf("%s%s %s", n.Serial, n.Area, n.Sequence)
}

func (n *NI1966) acceptableArea(area string) bool {
	for _, r := range area {
		if r == 'I' || r == 'Z' {
			return true
		}
	}

	return false
}
