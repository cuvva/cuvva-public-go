package cher

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type testErr struct {
	error   string
	timeout bool
}

func (e testErr) Error() string { return e.error }
func (e testErr) Timeout() bool { return e.timeout }

func TestCoerceThirdPartyTimeout(t1 *testing.T) {
	tests := []struct {
		name   string
		err    error
		expect func(t *testing.T, err error)
	}{
		{
			name: "wrapped timeout error",
			err: fmt.Errorf("wrap: %w", testErr{
				error:   "foobar",
				timeout: true,
			}),
			expect: func(t *testing.T, err error) {
				var cErr E
				require.ErrorAs(t, err, &cErr)
				require.Equal(t, "third_party_timeout", cErr.Code)
				require.Equal(t, "foobar", cErr.Meta["error"])
			},
		},
		{
			name: "other error",
			err:  fmt.Errorf("any error"),
			expect: func(t *testing.T, err error) {
				require.EqualError(t, err, "any error")
			},
		},
		{
			name: "no error",
			err:  nil,
			expect: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
	}
	for _, tc := range tests {
		t1.Run(tc.name, func(t *testing.T) {
			err := CoerceThirdPartyTimeout(tc.err)
			tc.expect(t, err)
		})
	}
}
