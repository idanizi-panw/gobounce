package gobounce

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestThrottle(t *testing.T) {
	t.Run("should return throttled function and cancel", func(t *testing.T) {
		f := func() {}
		throttled, cancel := Throttle(f, 0, nil)
		defer cancel()
		assert.IsType(t, func() {}, throttled)
		assert.IsType(t, func() {}, cancel)
		assert.NotNil(t, throttled)
		assert.NotNil(t, cancel)
	})

	t.Run("should invoke throttled function in wait time-span only once", func(t *testing.T) {
		var counter int
		f := func() { counter++ }
		wait := 1 * time.Second
		tick := 100 * time.Millisecond
		throttled, cancel := Throttle(f, wait, nil)
		defer cancel()

		require.Never(t, func() bool {
			for i := 0; i < 10; i++ {
				throttled()
			}
			return counter > 1
		}, wait+(2*tick), tick)
		assert.Equal(t, 1, counter)

		require.Never(t, func() bool {
			for i := 0; i < 10; i++ {
				throttled()
			}
			return counter > 2
		}, wait+(2*tick), tick)
		assert.Equal(t, 2, counter)
	})

	t.Run("should cancel throttling when cancel called", func(t *testing.T) {
		var counter int
		f := func() { counter++ }
		wait := 1 * time.Second
		tick := 100 * time.Millisecond
		throttled, cancel := Throttle(f, wait, nil)
		require.Never(t, func() bool {
			for i := 0; i < 10; i++ {
				throttled()
			}
			return counter > 1
		}, wait+(2*tick), tick)
		assert.Equal(t, 1, counter)

		cancel()
		require.Eventually(t, func() bool {
			return counter == 2
		}, 2*tick, tick)

		require.Never(t, func() bool {
			for i := 0; i < 10; i++ {
				throttled()
			}
			return counter > 2
		}, wait+(2*tick), tick)
		assert.Equal(t, 2, counter)
	})

	t.Run("should call once in time span with multiple goroutines", func(t *testing.T) {
		var counter int
		var i int
		wait := 2 * time.Second
		tick := 100 * time.Millisecond
		throttled, cancel := Throttle(func() {
			counter++
		}, wait, nil)
		defer cancel()
		for i = 0; i < 3; i++ {
			go func() {
				for j := 0; j < 10; j++ {
					throttled()
				}
			}()
		}
		require.Never(t, func() bool {
			return counter > 1
		}, wait+(2*tick), tick)
		assert.Equal(t, 1, counter)
	})

	t.Run("should call once in time span, cancel, with multiple goroutines", func(t *testing.T) {
		var counter int
		wait := 2 * time.Second
		tick := 100 * time.Millisecond
		throttled, cancel := Throttle(func() { counter++ }, wait, nil)

		require.Never(t, func() bool {
			go throttled()
			return counter > 1
		}, wait+(2*tick), tick)
		assert.Equal(t, 1, counter)

		cancel()
		require.Eventually(t, func() bool {
			return counter == 2
		}, 2*tick, tick)

		require.Never(t, func() bool {
			go throttled()
			return counter > 2
		}, wait+(2*tick), tick)
		assert.Equal(t, 2, counter)
	})

	t.Run("should call function on trailing edge by default", func(t *testing.T) {
		var counter int
		wait := 2 * time.Second
		tick := 100 * time.Millisecond
		throttled, cancel := Throttle(func() { counter++ }, wait, nil)
		defer cancel()

		require.Never(t, func() bool {
			go throttled()
			return counter > 0
		}, (wait/2)+(2*tick), tick)
		assert.Equal(t, 0, counter)

		require.Never(t, func() bool {
			go throttled()
			return counter > 1
		}, (wait/2)+(2*tick), tick)
		assert.Equal(t, 1, counter)
	})

	t.Run("should call function on leading edge", func(t *testing.T) {
		var counter int
		wait := 2 * time.Second
		tick := 100 * time.Millisecond
		options := &ThrottleOptions{
			Trailing: false,
			Leading:  true,
			Ctx:      context.Background(),
		}
		throttled, cancel := Throttle(func() { counter++ }, wait, options)
		defer cancel()

		require.Never(t, func() bool {
			go throttled()
			return counter > 1
		}, wait/2, tick)
		assert.Equal(t, 1, counter)

		require.Never(t, func() bool {
			go throttled()
			return counter > 1
		}, wait/2, tick)
		assert.Equal(t, 1, counter)
	})

	t.Run("should call function twice if trailing and leading are true", func(t *testing.T) {
		options := &ThrottleOptions{
			Trailing: true,
			Leading:  true,
			Ctx:      context.Background(),
		}
		var counter int
		wait := 1 * time.Second
		tick := 100 * time.Millisecond
		throttled, cancel := Throttle(func() { counter++ }, wait, options)
		defer cancel()

		require.Never(t, func() bool {
			go throttled()
			return counter > 1
		}, wait/2, tick)
		assert.Equal(t, 1, counter)

		require.Never(t, func() bool {
			go throttled()
			return counter > 2
		}, wait/2+tick, tick)
		assert.Equal(t, 2, counter)
	})

}
