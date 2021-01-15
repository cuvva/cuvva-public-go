package workerpool

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var inputs = []string{"foo", "bar", "wibble", "wobble", "wubble", "flob"} // this is the real list

func TestRunsProcesses(t *testing.T) {
	wp := New(context.Background(), 1)

	responseChannel := make(chan string, len(inputs))
	for _, i := range inputs {
		i := i
		wp.Do(func(ctx context.Context) error {
			responseChannel <- i
			return nil
		})
	}

	err := wp.Wait()
	close(responseChannel)

	assert.NoError(t, err)

	responseSet := channelToSet(responseChannel)
	assert.Equal(t, len(responseSet), len(inputs))
	for _, input := range inputs {
		assert.Contains(t, responseSet, input)
	}
}

func TestHandlesErrors(t *testing.T) {
	wp := New(context.Background(), 3)

	responseChannel := make(chan string, len(inputs))
	for _, input := range inputs {
		input := input
		wp.Do(func(ctx context.Context) error {
			if input == "wibble" {
				return errors.New("wibble error")
			}

			responseChannel <- input
			return nil
		})
	}

	err := wp.Wait()
	close(responseChannel)

	assert.Error(t, err)
	assert.Equal(t, "wibble error", err.Error())
	responseSet := channelToSet(responseChannel)
	assert.NotContains(t, responseSet, "wibble")
}

func TestExecutesImmediately(t *testing.T) {
	wp := New(context.Background(), 1)

	responseChannel := make(chan string, len(inputs))
	for _, input := range inputs {
		input := input
		wp.Do(func(ctx context.Context) error {
			responseChannel <- input
			return nil
		})
	}

Exit:
	for {
		select {
		case <-responseChannel:
			break Exit // received an item
		case <-time.After(10 * time.Second):
			t.Errorf("worker pool did not execute in time")
		}
	}

	wp.Wait()
	close(responseChannel)
}

func channelToSet(channel chan string) map[string]struct{} {
	responseSet := map[string]struct{}{}
	for item := range channel {
		responseSet[item] = struct{}{}
	}

	return responseSet
}
