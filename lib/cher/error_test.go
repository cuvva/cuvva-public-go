package cher

import (
	"errors"
	"net/http"
	"testing"

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
				E{Code: "foo"},
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
