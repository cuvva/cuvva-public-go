package ab

import (
	"context"
	"testing"

	"github.com/cuvva/cuvva-public-go/lib/ksuid"
)

func TestCheck(t *testing.T) {
	cases := []struct {
		test     string
		sub      string
		count    uint16
		expected bool
	}{
		{"test_1", "user_1", 20000, false},
		{"test_2", "user_1", 20000, false},
		{"test_3", "user_1", 20000, false},
		{"test_4", "user_1", 20000, false},
		{"test_5", "user_1", 20000, true},
		{"test_5", "user_1", 9931, false},
		{"test_5", "user_1", 9932, true},
		{"test_6", "user_1", 20000, true},
		{"test_1", "user_2", 20000, false},
		{"test_1", "user_3", 20000, true},
		{"test_1", "user_4", 20000, false},
		{"test_1", "user_5", 20000, true},
		{"test_1", "user_6", 20000, false},
	}

	for _, c := range cases {
		result := Check(c.test, c.sub, c.count)

		if result != c.expected {
			t.Errorf("%s, %s, %d - returned %t, expected %t", c.test, c.sub, c.count, result, c.expected)
		}
	}
}

var testID = ksuid.Generate(context.Background(), "test").String()
var userID = ksuid.Generate(context.Background(), "user").String()

func BenchmarkParse(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Check(testID, userID, 5000)
	}
}
