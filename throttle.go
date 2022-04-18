package gobounce

import (
	"context"
	"time"
)

//ThrottleOptions ...
type ThrottleOptions struct {
	Leading  bool
	Trailing bool // default true
	Ctx      context.Context
}

//NewThrottleOptions by default Trailing is true, Leading is false, and Ctx is context.Background()
func NewThrottleOptions() *ThrottleOptions {
	return &ThrottleOptions{Leading: false, Trailing: true, Ctx: context.Background()}
}

//Throttle Creates a throttled function that only invokes func at most once per every
// wait milliseconds. The throttled function comes with a cancel method to cancel
// delayed func invocations and a flush method to immediately invoke them.
// Provide options to indicate whether func should be invoked on the leading
// and/or trailing edge of the wait timeout. The func is invoked with the last
// arguments provided to the throttled function. Subsequent calls to the
// throttled function return the result of the last func invocation.
//
// Inspired by Lodash Throttle.
// See: https://lodash.com/docs/4.17.15#throttle
func Throttle(f func(), wait time.Duration, options *ThrottleOptions) (throttled func(), cancel func()) {
	if options == nil {
		options = NewThrottleOptions()
	}

	options.Ctx, cancel = context.WithCancel(options.Ctx)
	invoke := make(chan interface{})
	var lastInvoked time.Time
	var invokedCount uint64
	timer := time.NewTimer(wait)

	go func() {
		defer timer.Stop()
		for {
			select {
			case <-options.Ctx.Done():
				if invokedCount > 0 {
					go f()
				}
				return
			case <-invoke:
				invokedCount++
				if !lastInvoked.IsZero() && time.Since(lastInvoked) < wait {
					break
				}
				lastInvoked = time.Now()
				if options.Leading {
					go f()
				}
			case <-timer.C:
				timer.Reset(wait)

				if invokedCount == 0 {
					break
				}

				if !options.Trailing {
					invokedCount = 0
					break
				}

				if options.Leading {
					if invokedCount > 1 {
						invokedCount = 0
						go f()
						break
					}
					break
				}

				invokedCount = 0
				go f()
			}
		}
	}()

	throttled = func() {
		invoke <- struct{}{}
	}

	return throttled, cancel
}
