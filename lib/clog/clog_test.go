package clog

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestContextLogger(t *testing.T) {
	t.Run("Set", func(t *testing.T) {
		log := logrus.New().WithField("foo", "bar")

		r := &http.Request{}
		r = r.WithContext(Set(r.Context(), log))

		l := Get(r.Context())

		assert.Equal(t, log, l)
	})

	t.Run("SetFields", func(t *testing.T) {
		log := logrus.New().WithField("foo", "bar")

		r := &http.Request{}
		r = r.WithContext(Set(r.Context(), log))

		err := SetFields(r.Context(), Fields{
			"foo2": "bar2",
		})
		if err != nil {
			t.Fatal(err)
		}

		cl := getContextLogger(r.Context()).GetLogger()
		assert.Equal(t, "bar", cl.Data["foo"])
		assert.Equal(t, "bar2", cl.Data["foo2"])
	})

	t.Run("SetField", func(t *testing.T) {
		log := logrus.New().WithField("foo", "bar")

		r := &http.Request{}
		r = r.WithContext(Set(r.Context(), log))

		err := SetField(r.Context(), "foo2", "bar2")
		if err != nil {
			t.Fatal(err)
		}

		cl := getContextLogger(r.Context()).GetLogger()
		assert.Equal(t, "bar", cl.Data["foo"])
		assert.Equal(t, "bar2", cl.Data["foo2"])
	})

	t.Run("SetError", func(t *testing.T) {
		log := logrus.New().WithField("foo", "bar")

		r := &http.Request{}
		r = r.WithContext(Set(r.Context(), log))

		testError := errors.New("test error")

		err := SetError(r.Context(), testError)
		if err != nil {
			t.Fatal(err)
		}

		cl := getContextLogger(r.Context()).GetLogger()
		assert.Equal(t, testError, cl.Data["error"])
	})

	t.Run("Logger when no clog is set", func(t *testing.T) {
		r := &http.Request{}
		l := Get(r.Context())

		assert.NotNil(t, l)
		assert.IsType(t, &logrus.Entry{}, l)
	})

	t.Run("SetField when no logger is set", func(t *testing.T) {
		err := SetField(context.Background(), "foo", "bar")

		assert.NotNil(t, err)
		assert.Equal(t, "no clog exists in the context", err.Error())
	})

	t.Run("SetFields when no logger is set", func(t *testing.T) {
		err := SetFields(context.Background(), Fields{"foo": "bar"})

		assert.NotNil(t, err)
		assert.Equal(t, "no clog exists in the context", err.Error())
	})

	t.Run("SetError when no logger is set", func(t *testing.T) {
		err := SetError(context.Background(), errors.New("foo"))

		assert.NotNil(t, err)
		assert.Equal(t, "no clog exists in the context", err.Error())
	})

}
