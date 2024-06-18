package retrier

import (
	"context"
	"math"
	"time"
)

type Strategy func(ctx context.Context) chan struct{}

func Backoff(repeats int, duration time.Duration, factor float64, maxDelay time.Duration) Strategy {
	return func(ctx context.Context) chan struct{} {
		ch := make(chan struct{})
		go func() {
			defer close(ch)

			for i := 0; i < repeats; i++ {
				select {
				case <-ctx.Done():
					return
				case ch <- struct{}{}:
				}

				delay := time.Duration(
					math.Min(
						float64(maxDelay),
						float64(duration)*math.Pow(factor, float64(i)),
					),
				)
				sleep(ctx, delay)
			}
		}()
		return ch
	}
}

func sleep(ctx context.Context, duration time.Duration) {
	select {
	case <-time.After(duration):
		return
	case <-ctx.Done():
		return
	}
}
