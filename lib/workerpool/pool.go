package workerpool

import (
	"context"
	"sync/atomic"
)

type Worker struct {
	ctx              context.Context
	cancel           context.CancelFunc
	err              *errorContainer
	pendingTasks     chan WorkerFunc
	activeOperations int64
}

type WorkerFunc func(context.Context) error

func New(ctx context.Context, workers int) *Worker {
	cCtx, ctxCancel := context.WithCancel(ctx)

	wrk := &Worker{
		ctx:              cCtx,
		cancel:           ctxCancel,
		err:              &errorContainer{},
		pendingTasks:     make(chan WorkerFunc),
		activeOperations: 0,
	}
	wrk.start(workers)
	return wrk
}

func (w *Worker) start(numWorkers int) {
	for i := 0; i < numWorkers; i++ {
		go func() {
			for {
				select {
				case fn, ok := <-w.pendingTasks:
					if !ok {
						return
					}

					if err := fn(w.ctx); err != nil {
						w.err.AssignError(err)
					}

					atomic.AddInt64(&w.activeOperations, -1)
				case <-w.ctx.Done():
					if w.ctx.Err() != nil {
						w.err.AssignError(w.ctx.Err())
					}
					return
				default:
					if w.err.err != nil {
						w.cancel()
					}
				}
			}
		}()
	}
}

func (w *Worker) Do(fns ...WorkerFunc) {
	for _, fn := range fns {
		select {
		case <-w.ctx.Done():
			return
		default:
			atomic.AddInt64(&w.activeOperations, 1)
			f := fn
			go func() {
				defer func() {
					// If we couldn't send (channel closed), decrement the counter
					if r := recover(); r != nil {
						atomic.AddInt64(&w.activeOperations, -1)
					}
				}()
				// Check context first before attempting to send
				select {
				case <-w.ctx.Done():
					atomic.AddInt64(&w.activeOperations, -1)
					return
				default:
					// Try to send, but handle panic if channel is closed
					func() {
						defer func() {
							if r := recover(); r != nil {
								atomic.AddInt64(&w.activeOperations, -1)
							}
						}()
						select {
						case <-w.ctx.Done():
							atomic.AddInt64(&w.activeOperations, -1)
							return
						case w.pendingTasks <- f:
							// Successfully sent
						}
					}()
				}
			}()
		}
	}
}

func (w *Worker) Wait() error {
	defer close(w.pendingTasks)

	for {
		select {
		case <-w.ctx.Done():
			return w.err.err
		default:
			if atomic.LoadInt64(&w.activeOperations) == 0 {
				return w.err.err
			}
		}
	}
}
