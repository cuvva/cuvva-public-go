package postcode

import "testing"

func TestValidation(t *testing.T) {
	for input, expected := range expectedValidation {
		postcode, err := Parse(input)

		if expected {
			if postcode == nil || err != nil {
				t.Errorf("should have validated - %s", input)
			}
		} else {
			if postcode != nil || err == nil {
				t.Errorf("should have failed - %s", input)
			}
		}
	}
}

func TestNormalized(t *testing.T) {
	for input, expected := range expectedNormalized {
		postcode, err := Parse(input)
		if err != nil {
			t.Errorf("unexpected error - %s", input)
		}

		actual := postcode.FullNormalized
		if expected != actual {
			t.Errorf("unexpected result - expected '%s', got '%s'", expected, actual)
		}
	}
}

func TestCompact(t *testing.T) {
	for input, expected := range expectedCompact {
		postcode, err := Parse(input)
		if err != nil {
			t.Errorf("unexpected error - %s", input)
		}

		actual := postcode.FullCompact
		if expected != actual {
			t.Errorf("unexpected result - expected '%s', got '%s'", expected, actual)
		}
	}
}

func TestArea(t *testing.T) {
	for input, expected := range expectedArea {
		postcode, err := Parse(input)
		if err != nil {
			t.Errorf("unexpected error - %s", input)
		}

		actual := postcode.Area
		if expected != actual {
			t.Errorf("unexpected result - expected '%s', got '%s'", expected, actual)
		}
	}
}

func TestDistrict(t *testing.T) {
	for input, expected := range expectedDistrict {
		postcode, err := Parse(input)
		if err != nil {
			t.Errorf("unexpected error - %s", input)
		}

		actual := postcode.District
		if expected != actual {
			t.Errorf("unexpected result - expected '%s', got '%s'", expected, actual)
		}
	}
}

func TestSubDistrict(t *testing.T) {
	for input, expected := range expectedSubDistrict {
		postcode, err := Parse(input)
		if err != nil {
			t.Errorf("unexpected error - %s", input)
		}

		actual := postcode.SubDistrict
		if expected != actual {
			t.Errorf("unexpected result - expected '%s', got '%s'", expected, actual)
		}
	}
}

func TestOutcodes(t *testing.T) {
	for input, expected := range expectedOutcode {
		postcode, err := Parse(input)
		if err != nil {
			t.Errorf("unexpected error - %s", input)
		}

		actual := postcode.Outcode
		if expected != actual {
			t.Errorf("unexpected result - expected '%s', got '%s'", expected, actual)
		}
	}
}

func TestSector(t *testing.T) {
	for input, expected := range expectedSector {
		postcode, err := Parse(input)
		if err != nil {
			t.Errorf("unexpected error - %s", input)
		}

		actual := postcode.Sector
		if expected != actual {
			t.Errorf("unexpected result - expected '%s', got '%s'", expected, actual)
		}
	}
}

func TestIncodes(t *testing.T) {
	for input, expected := range expectedIncode {
		postcode, err := Parse(input)
		if err != nil {
			t.Errorf("unexpected error - %s", input)
		}

		actual := postcode.Incode
		if expected != actual {
			t.Errorf("unexpected result - expected '%s', got '%s'", expected, actual)
		}
	}
}

func TestUnit(t *testing.T) {
	for input, expected := range expectedUnit {
		postcode, err := Parse(input)
		if err != nil {
			t.Errorf("unexpected error - %s", input)
		}

		actual := postcode.Unit
		if expected != actual {
			t.Errorf("unexpected result - expected '%s', got '%s'", expected, actual)
		}
	}
}
