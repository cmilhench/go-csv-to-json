package pipeline

import "context"

// Map applies a function to values received from an input channel.
// It sends the results of the function to an output channel.
// The function stops when the input channel is closed or the context is canceled.
func Map[I, O any](ctx context.Context, in <-chan I, fn func(context.Context, int, I) O) <-chan O {
	out := make(chan O)
	go func() {
		defer close(out)
		for i := 0; ; i++ {
			select {
			case <-ctx.Done():
				return
			case v, ok := <-in:
				if !ok {
					return
				}
				select {
				case out <- fn(ctx, i, v):
				case <-ctx.Done():
					return
				}
			}
		}
	}()
	return out
}

// Take receives values from an input channel and sends only the first n values
// to the output channel.
func Take[T any](ctx context.Context, in <-chan T, n int) <-chan T {
	out := make(chan T)
	go func() {
		defer close(out)
		for i := 0; i < n; i++ {
			select {
			case <-ctx.Done():
				return
			case out <- <-in: // out ← in
			}
		}
	}()
	return out
}
