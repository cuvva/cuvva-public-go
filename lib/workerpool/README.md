# workerpool

Used to execute tasks with a concurrent limit. The worker pool starts executing as soon as work is submitted.
If an error occurs during a piece of work, the worker pool stops executing all work and returns an error.

## Example

Here a worker pool executes max 2 requests concurrently.

```go
func (a *App) getIDs(ctx context.Context, ids []string) error {
	wp := workerpool.New(ctx, 2)

	for _, id := range ids {
		id := id
		wp.Do(func(ctx context.Context) {
			resp, err := a.ServiceClient.GetIDs(ctx, id)
			if err != nil {
				return err
			}

			// store response in concurrency safe store
		})
	}

	return wp.Wait()
}
```

