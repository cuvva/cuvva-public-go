package pg

import (
	"fmt"
	"testing"
)

func TestScan(t *testing.T) {
	p := Point{
		X: 53.4504009999999994,
		Y: -2.31376399999999993,
	}

	v, err := p.Value()
	if err != nil {
		t.Error(err)
	}

	b, ok := v.([]byte)
	if !ok {
		t.Errorf("incorrect type for driver value (expected []byte)")
	}

	str := string(b)
	expected := fmt.Sprintf(pointFormat, p.X, p.Y)

	if str != expected {
		t.Errorf("expected %s, got %s", expected, str)
	}
}

func TestValue(t *testing.T) {
	src := []byte("(53.450400999999999,-2.313764000000000)")

	var p Point

	if err := p.Scan(src); err != nil {
		t.Error(err)
	}

	if p.X != 53.450400999999999 {
		t.Errorf("incorrect X coord in value")
	}

	if p.Y != -2.313764000000000 {
		t.Errorf("incorrect Y coord in value")
	}
}
