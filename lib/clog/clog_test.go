package clog

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/cuvva/cuvva-public-go/lib/cher"
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

func TestDetermineLevel(t *testing.T) {
	type testCase struct {
		name             string
		err              error
		timeoutsAsErrors bool
		expected         logrus.Level
	}

	tests := []testCase{
		{
			name:             "bad request",
			err:              cher.New("bad_request", nil),
			timeoutsAsErrors: false,
			expected:         logrus.WarnLevel,
		},
		{
			name:             "context cancelled",
			err:              cher.New(cher.ContextCanceled, nil),
			timeoutsAsErrors: false,
			expected:         logrus.InfoLevel,
		},
		{
			name:             "context cancelled with timeouts as errors",
			err:              cher.New(cher.ContextCanceled, nil),
			timeoutsAsErrors: true,
			expected:         logrus.ErrorLevel,
		},
		{
			name:             "unknown",
			err:              cher.New(cher.Unknown, nil),
			timeoutsAsErrors: false,
			expected:         logrus.ErrorLevel,
		},
		{
			name:             "postgres context cancelled",
			err:              fmt.Errorf("pq: canceling statement due to user request"),
			timeoutsAsErrors: false,
			expected:         logrus.InfoLevel,
		},
		{
			name:             "postgres context cancelled with timeouts as errors",
			err:              fmt.Errorf("pq: canceling statement due to user request"),
			timeoutsAsErrors: true,
			expected:         logrus.ErrorLevel,
		},
		{
			name:             "other error",
			err:              fmt.Errorf("something something darkside"),
			timeoutsAsErrors: false,
			expected:         logrus.ErrorLevel,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := DetermineLevel(tc.err, tc.timeoutsAsErrors)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestWithCherFields(t *testing.T) {
	t.Run("With cher.E containing reasons and meta", func(t *testing.T) {
		entry := logrus.NewEntry(logrus.New())
		err := cher.E{
			Code:    "test_error",
			Reasons: []cher.E{{Code: "reason_1"}, {Code: "reason_2"}},
			Meta:    map[string]interface{}{"key": "value"},
		}

		result := WithCherFields(entry, err)

		assert.Equal(t, []cher.E{{Code: "reason_1"}, {Code: "reason_2"}}, result.Data["error_reasons"])
		assert.Equal(t, map[string]interface{}{"key": "value"}, result.Data["error_meta"])
	})

	t.Run("With cher.E containing no reasons or meta", func(t *testing.T) {
		entry := logrus.NewEntry(logrus.New())
		err := cher.E{
			Code: "test_error",
		}

		result := WithCherFields(entry, err)

		assert.NotContains(t, result.Data, "error_reasons")
		assert.NotContains(t, result.Data, "error_meta")
	})

	t.Run("With non-cher error", func(t *testing.T) {
		entry := logrus.NewEntry(logrus.New())
		err := errors.New("non-cher error")

		result := WithCherFields(entry, err)

		assert.NotContains(t, result.Data, "error_reasons")
		assert.NotContains(t, result.Data, "error_meta")
	})
}
