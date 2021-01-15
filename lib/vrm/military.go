package vrm

import (
	"strings"
)

type Military struct {
	Sections []string
}

func ParseMilitary(vrm string) VRM {
	if len(vrm) != 6 {
		return nil
	}

	if !(match(vrm[:2], isAlpha) && match(vrm[2:4], isNumeric) && match(vrm[4:], isAlpha)) &&
		!(match(vrm[:2], isNumeric) && match(vrm[2:4], isAlpha) && match(vrm[4:], isNumeric)) {
		return nil
	}

	return &Military{
		Sections: []string{
			vrm[:2], vrm[2:4], vrm[4:],
		},
	}
}

func (m *Military) Format() string {
	return "military"
}

func (m *Military) String() string {
	return strings.Join(m.Sections, "")
}

func (m *Military) PrettyString() string {
	return strings.Join(m.Sections, " ")
}
