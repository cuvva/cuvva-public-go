package request

import (
	"context"
	"errors"
	"testing"
	"time"
)

type testingContextKey string

var tcKey testingContextKey = "foo"

func Test_cloneContext(t *testing.T) {

	t.Run("values accessible", func(t *testing.T) {
		want := "bar"
		orig := context.WithValue(context.Background(), tcKey, want)
		forked := cloneContext(orig)

		if got := forked.Value(tcKey); got != want {
			t.Errorf(".Value() = %v, want %v", got, want)
		}
	})

	t.Run("modify forked without affecting original", func(t *testing.T) {
		orig := context.WithValue(context.Background(), tcKey, "bar1")
		forked := cloneContext(orig)
		forked = context.WithValue(forked, tcKey, "bar2")

		if got := orig.Value(tcKey); got != "bar1" {
			t.Errorf("orig.Value() = %v, want %v", got, "bar1")
		}
		if got := forked.Value(tcKey); got != "bar2" {
			t.Errorf("forked.Value() = %v, want %v", got, "bar2")
		}
	})

	t.Run("modify original without affecting forked", func(t *testing.T) {
		orig := context.WithValue(context.Background(), tcKey, "bar1")
		forked := cloneContext(orig)
		orig = context.WithValue(orig, tcKey, "bar2")

		if got := forked.Value(tcKey); got != "bar1" {
			t.Errorf("forked.Value() = %v, want %v", got, "bar1")
		}
		if got := orig.Value(tcKey); got != "bar2" {
			t.Errorf("orig.Value() = %v, want %v", got, "bar2")
		}
	})

	t.Run("removes cancellation", func(t *testing.T) {
		orig, cancelFn := context.WithCancel(context.Background())
		forked := cloneContext(orig)

		cancelFn()

		if !errors.Is(orig.Err(), context.Canceled) {
			t.Errorf("expected original context to have context canceled error")
		}
		if errors.Is(forked.Err(), context.Canceled) {
			t.Errorf("do not expect forked context to have context canceled error")
		}
	})

	t.Run("removes deadline", func(t *testing.T) {
		orig, cancelFn := context.WithTimeout(context.Background(), 1*time.Millisecond)
		defer cancelFn()
		forked := cloneContext(orig)

		time.Sleep(5 * time.Millisecond)

		if !errors.Is(orig.Err(), context.DeadlineExceeded) {
			t.Errorf("expected original context to have deadline exceeded error")
		}
		if errors.Is(forked.Err(), context.DeadlineExceeded) {
			t.Errorf("do not expect forked context to have deadline exceeded error")
		}
	})

}
