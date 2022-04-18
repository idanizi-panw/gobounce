package gobounce

import (
	"context"
	"math"
	"time"
)

//DebounceOptions ...
type DebounceOptions struct {
	ThrottleOptions
	MaxWait time.Duration
}

//NewDebounceOptions ...
func NewDebounceOptions() *DebounceOptions {
	return &DebounceOptions{
		ThrottleOptions: ThrottleOptions{Leading: false, Trailing: true, Ctx: context.Background()},
		MaxWait:         math.MaxInt64,
	}
}

//Debounce Creates a debounced function that delays invoking func until after wait
// milliseconds have elapsed since the last time the debounced function was
// invoked. The debounced function comes with a cancel method to cancel delayed
// func invocations and a flush method to immediately invoke them. Provide
// options to indicate whether func should be invoked on the leading and/or
// trailing edge of the wait timeout. The func is invoked with the last arguments
// provided to the debounced function. Subsequent calls to the debounced function
// return the result of the last func invocation.
//
// Inspired by Lodash Debounce.
// See: https://lodash.com/docs/4.17.15#debounce
func Debounce(f func(), wait time.Duration, options *DebounceOptions) (debounced func(), cancel func()) {
	if options == nil {
		options = NewDebounceOptions()
	}

	options.Ctx, cancel = context.WithCancel(options.Ctx)
	invoke := make(chan interface{})
	var lastInvoked time.Time
	var invokedCount uint64
	timer := time.NewTimer(wait)
	maxTimer := time.NewTimer(options.MaxWait)

	go func() {
		defer timer.Stop()
		defer maxTimer.Stop()
		for {
			select {
			case <-options.Ctx.Done():
				if invokedCount > 0 {
					go f()
				}
				return
			case <-maxTimer.C:
				if invokedCount > 0 {
					go f()
				}
				maxTimer.Reset(options.MaxWait)
			case <-invoke:
				invokedCount++
				timer.Reset(wait)
				if !lastInvoked.IsZero() && time.Since(lastInvoked) < wait {
					break
				}
				lastInvoked = time.Now()
				if options.Leading {
					maxTimer.Reset(options.MaxWait)
					go f()
				}
			case <-timer.C:
				if !options.Trailing {
					invokedCount = 0
					break
				}

				if options.Leading {
					if invokedCount > 1 {
						go f()
						break
					}
					break
				}

				if invokedCount > 0 {
					go f()
				}
			}
		}
	}()

	debounced = func() {
		invoke <- struct{}{}
	}

	return debounced, cancel
}
