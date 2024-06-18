package retrier

import (
	"context"
	"errors"
)

type Retrier struct {
	Strategy Strategy
	OnRetry  func(ctx context.Context, n int, err error)
}

func (r *Retrier) Do(ctx context.Context, fn func() error, errs ...error) (err error) {
	ctx, cancelFunc := context.WithCancel(ctx)
	defer cancelFunc()

	attempt := 0
	ch := r.Strategy(ctx)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case _, ok := <-ch:
			if !ok {
				return err
			}

			if r.OnRetry != nil && attempt > 0 {
				r.OnRetry(ctx, attempt, err)
			}

			if err = fn(); err == nil {
				return nil
			}
			if len(errs) > 0 && !oneOfErrs(err, errs...) {
				return err
			}
			attempt++
		}
	}
}

func oneOfErrs(err error, errs ...error) bool {
	for _, e := range errs {
		if errors.Is(err, e) {
			return true
		}
	}
	return false
}
