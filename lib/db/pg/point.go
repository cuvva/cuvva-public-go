package pg

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// Point is the representation of postgres type 'point' (x,y)
type Point struct {
	X, Y float64
}

const pointFormat = "(%.15f,%.15f)"

// String implements the stringer interface and converts the point into the
// correct string display type
func (p Point) String() string {
	return fmt.Sprintf("POINT"+pointFormat, p.X, p.Y)
}

// Value implements driver.Valuer and coerces Point types into the native
// postgres 'point' type. It retains 15 digits of precision, as required for
// float8/double precision.
func (p Point) Value() (driver.Value, error) {
	return []byte(fmt.Sprintf(pointFormat, p.X, p.Y)), nil
}

// Scan implements sql.Scanner and sets the X/Y coordinates of the point with
// the given src
func (p *Point) Scan(src interface{}) error {
	raw, ok := src.([]uint8)
	if !ok {
		return errors.New("pg: point type not valid")
	}

	coords := strings.Split(string(raw), ",")
	if len(coords) != 2 {
		return errors.New("incorrect parts for point")
	}

	for i, part := range coords {
		coords[i] = strings.TrimFunc(part, func(r rune) bool {
			if r == '(' || r == ')' {
				return true
			}

			return false
		})
	}

	x, err := strconv.ParseFloat(coords[0], 64)
	if err != nil {
		return err
	}

	y, err := strconv.ParseFloat(coords[1], 64)
	if err != nil {
		return err
	}

	*p = Point{
		X: x,
		Y: y,
	}

	return nil
}
