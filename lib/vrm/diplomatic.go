package vrm

import (
	"fmt"
	"strconv"
)

type Diplomatic struct {
	Type byte

	Entity int
	Serial int
}

func ParseDiplomatic(vrm string) VRM {
	if len(vrm) != 7 {
		return nil
	}

	if vrm[3] != 'D' && vrm[3] != 'X' {
		return nil
	}

	if !match(vrm[:3], isNumeric) {
		return nil
	} else if !match(vrm[4:], isNumeric) {
		return nil
	}

	entity, _ := strconv.Atoi(vrm[:3])
	serial, _ := strconv.Atoi(vrm[4:])

	return &Diplomatic{
		Type: vrm[3],

		Entity: entity,
		Serial: serial,
	}
}

func (d *Diplomatic) Format() string {
	return "diplomatic"
}

func (d *Diplomatic) String() string {
	return fmt.Sprintf("%03d%c%03d", d.Entity, d.Type, d.Serial)
}

func (d *Diplomatic) PrettyString() string {
	return fmt.Sprintf("%03d %c %03d", d.Entity, d.Type, d.Serial)
}
