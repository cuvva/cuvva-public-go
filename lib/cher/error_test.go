package cher

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestE(t *testing.T) {
	t.Run("New", func(t *testing.T) {
		m := M{"foo": "bar"}
		e := New(NotFound, m, E{Code: "foo"})

		assert.Equal(t, e, E{
			Code: NotFound,
			Meta: m,
			Reasons: []E{
				{Code: "foo"},
			},
		})
	})

	t.Run("Errorf", func(t *testing.T) {
		m := M{"foo": "bar"}
		e := Errorf(NotFound, m, "foo %s", "bar")

		assert.Equal(t, e, E{
			Code: NotFound,
			Meta: M{
				"foo":     "bar",
				"message": "foo bar",
			},
		})
	})

	t.Run("StatusCode", func(t *testing.T) {
		tests := []struct {
			Name       string
			E          E
			StatusCode int
		}{
			{"BadRequest", E{Code: BadRequest}, http.StatusBadRequest},
			{"Unauthorized", E{Code: Unauthorized}, http.StatusUnauthorized},
			{"AccessDenied", E{Code: AccessDenied}, http.StatusForbidden},
			{"NotFound", E{Code: NotFound}, http.StatusNotFound},
			{"Unknown", E{Code: Unknown}, http.StatusInternalServerError},
			{"Handled", E{Code: "some_developer_code"}, http.StatusBadRequest},
		}

		for _, test := range tests {
			t.Run(test.Name, func(t *testing.T) {
				sc := test.E.StatusCode()
				assert.Equal(t, test.StatusCode, sc)
			})
		}
	})

	t.Run("Error", func(t *testing.T) {
		e := E{Code: NotFound}
		assert.Equal(t, NotFound, e.Error())
	})
}

func TestCoerce(t *testing.T) {
	tests := []struct {
		Name   string
		Src    interface{}
		Result E
	}{
		{"E", E{Code: "foo"}, E{Code: "foo"}},
		{"String", "foo", E{Code: "foo"}},
		{"JSON", []byte(`{"code":"foo"}`), E{Code: "foo"}},
		{"BadJSON", []byte(`{"code":0}`), E{Code: CoercionError, Meta: M{"message": "json: cannot unmarshal number into Go struct field E.code of type string"}}},
		{"Error", errors.New("foo"), E{Code: Unknown, Meta: M{"message": "foo"}}},
		{"Unknown", nil, E{Code: CoercionError}},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			e := Coerce(test.Src)

			assert.Equal(t, test.Result, e)
		})
	}
}

func TestWrapIfNotCher(t *testing.T) {
	type testCase struct {
		name   string
		msg    string
		err    error
		expect func(*testing.T, error)
	}

	tests := []testCase{
		{
			name: "nil",
			msg:  "foo",
			err:  nil,
			expect: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name: "err",
			msg:  "foo",
			err:  fmt.Errorf("nope"),
			expect: func(t *testing.T, err error) {
				assert.EqualError(t, err, "foo: nope")
			},
		},
		{
			name: "cher",
			msg:  "foo",
			err:  New("nope", nil),
			expect: func(t *testing.T, err error) {
				cErr, ok := err.(E)
				assert.True(t, ok)
				assert.Equal(t, "nope", cErr.Code)
			},
		},
		{
			name: "cher unknown",
			msg:  "foo",
			err:  New("unknown", nil),
			expect: func(t *testing.T, err error) {
				assert.EqualError(t, err, "foo: unknown")
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := WrapIfNotCher(tc.err, tc.msg)
			tc.expect(t, result)
		})
	}
}

func TestWrapIfNotCherCodes(t *testing.T) {
	type testCase struct {
		name   string
		msg    string
		err    error
		expect func(*testing.T, error)
	}

	tests := []testCase{
		{
			name: "nil",
			msg:  "foo",
			err:  nil,
			expect: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name: "err",
			msg:  "foo",
			err:  fmt.Errorf("nope"),
			expect: func(t *testing.T, err error) {
				assert.EqualError(t, err, "foo: nope")
			},
		},
		{
			name: "cher specified code",
			msg:  "foo",
			err:  New("code_1", nil),
			expect: func(t *testing.T, err error) {
				cErr, ok := err.(E)
				assert.True(t, ok)
				assert.Equal(t, "code_1", cErr.Code)
			},
		},
		{
			name: "cher other code",
			msg:  "foo",
			err:  New("unknown", nil),
			expect: func(t *testing.T, err error) {
				assert.EqualError(t, err, "foo: unknown")
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := WrapIfNotCher(tc.err, tc.msg)
			tc.expect(t, result)
		})
	}
}

func TestIsCherWithCode(t *testing.T) {
	type testCase struct {
		name   string
		err    error
		codes  []string
		expect func(t *testing.T, cErr E, ok bool)
	}

	tests := []testCase{
		{
			name:  "nil",
			err:   nil,
			codes: []string{"code_1"},
			expect: func(t *testing.T, cErr E, ok bool) {
				assert.False(t, ok)
			},
		},
		{
			name:  "normal error",
			err:   fmt.Errorf("nope"),
			codes: []string{"code_1"},
			expect: func(t *testing.T, cErr E, ok bool) {
				assert.False(t, ok)
			},
		},
		{
			name:  "normal error with same string",
			err:   fmt.Errorf("code_1"),
			codes: []string{"code_1"},
			expect: func(t *testing.T, cErr E, ok bool) {
				assert.False(t, ok)
			},
		},
		{
			name:  "cher specified code",
			err:   New("code_1", nil),
			codes: []string{"code_1"},
			expect: func(t *testing.T, cErr E, ok bool) {
				assert.True(t, ok)
				assert.Equal(t, "code_1", cErr.Code)
			},
		},
		{
			name:  "cher other code",
			err:   New("unknown", nil),
			codes: []string{"code_1"},
			expect: func(t *testing.T, cErr E, ok bool) {
				assert.False(t, ok)
			},
		},
		{
			name:  "wrapped cher",
			err:   errors.Wrap(New("code_1", nil), "wrapped"),
			codes: []string{"code_1"},
			expect: func(t *testing.T, cErr E, ok bool) {
				assert.True(t, ok)
				assert.Equal(t, "code_1", cErr.Code)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cErr, ok := IsCherWithCode(tc.err, tc.codes...)
			tc.expect(t, cErr, ok)
		})
	}
}
