package vrm

import (
	"fmt"
	"strconv"
)

type GB2001 struct {
	Area string

	FirstHalf bool

	Year int

	Serial string
}

func ParseGB2001(vrm string) VRM {
	if len(vrm) != 7 {
		return nil
	}

	g := &GB2001{
		Area:   vrm[:2],
		Serial: vrm[4:],
	}

	ageID := vrm[2:4]

	if !match(g.Area, isAlpha) || any(g.Area, g.isProhibitedLetterArea) {
		return nil
	} else if !match(ageID, isNumeric) || ageID == "01" {
		return nil
	} else if !match(g.Serial, isAlpha) || any(g.Serial, g.isProhibitedLetterSerial) {
		return nil
	}

	age, _ := strconv.Atoi(ageID) // already validated as numeric, can ignore error
	g.FirstHalf, g.Year = g.calcYear(age)

	return g
}

func (g *GB2001) Format() string {
	return "gb_2001"
}

func (g *GB2001) String() string {
	return fmt.Sprintf("%s%02d%s", g.Area, g.calcAge(g.FirstHalf, g.Year), g.Serial)
}

func (g *GB2001) PrettyString() string {
	return fmt.Sprintf("%s%02d %s", g.Area, g.calcAge(g.FirstHalf, g.Year), g.Serial)
}

func (g *GB2001) isProhibitedLetterArea(r rune) bool {
	return r == 'I' || r == 'Q' || r == 'Z'
}

func (g *GB2001) isProhibitedLetterSerial(r rune) bool {
	return r == 'I' || r == 'Q'
}

func (g *GB2001) calcYear(age int) (firstHalf bool, year int) {
	if age == 0 {
		return false, 2050
	} else if age > 50 {
		return false, (age - 50) + 2000
	}

	return true, age + 2000
}

func (g *GB2001) calcAge(firstHalf bool, year int) int {
	if !firstHalf {
		return (year + 50 - 2000) % 100
	}

	return year - 2000
}
