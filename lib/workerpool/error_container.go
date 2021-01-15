package workerpool

import (
	"sync"
)

type errorContainer struct {
	sync.Once
	err error
}

func (ec *errorContainer) AssignError(err error) {
	ec.Do(func() {
		ec.err = err
	})
}
