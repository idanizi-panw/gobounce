package gobounce

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDebounce(t *testing.T) {
	t.Run("should return debounce function and cancel", func(t *testing.T) {
		var counter int
		f := func() { counter++ }
		wait := 1 * time.Second
		debounced, cancel := Debounce(f, wait, nil)
		defer cancel()
		require.NotNil(t, debounced)
		require.NotNil(t, cancel)
		require.IsType(t, func() {}, debounced)
		require.IsType(t, func() {}, cancel)
	})

	t.Run("should call debounce function only after wait time-span had passed since last invocation", func(t *testing.T) {
		var counter int
		f := func() { counter++ }
		wait := 1 * time.Second
		tick := 100 * time.Millisecond
		debounced, cancel := Debounce(f, wait, nil)
		defer cancel()

		for i := 0; i < 10; i++ {
			debounced()
		}

		require.Never(t, func() bool {
			return counter > 0
		}, wait-tick, tick)
		require.Eventually(t, func() bool {
			return counter > 0
		}, 3*tick, tick)
		assert.Equal(t, 1, counter)
	})

	t.Run("should cancel debounce when cancel called", func(t *testing.T) {
		var counter int
		f := func() { counter++ }
		wait := 1 * time.Second
		tick := 100 * time.Millisecond
		debounced, cancel := Debounce(f, wait, nil)

		for i := 0; i < 15; i++ {
			debounced()
		}
		require.Never(t, func() bool {
			return counter > 0
		}, wait-tick, tick)

		for i := 0; i < 15; i++ {
			debounced()
		}

		require.Never(t, func() bool {
			return counter > 0
		}, wait-tick, tick)

		require.Eventually(t, func() bool {
			return counter == 1
		}, 3*tick, tick)

		cancel()
		require.Eventually(t, func() bool {
			return counter == 2
		}, 2*tick, tick)
	})

	t.Run("should invoked twice if leading and trailing are true and invoked more than once during time span", func(t *testing.T) {
		var counter int
		f := func() { counter++ }
		wait := 1 * time.Second
		tick := 100 * time.Millisecond
		options := NewDebounceOptions()
		options.Leading = true
		options.Trailing = true
		debounced, cancel := Debounce(f, wait, options)
		defer cancel()

		for i := 0; i < 10; i++ {
			debounced()
		}
		require.Never(t, func() bool {
			return counter > 1
		}, wait-tick, tick)
		require.Eventually(t, func() bool {
			return counter == 2
		}, 3*tick, tick)
		require.Never(t, func() bool {
			return counter > 2
		}, 2*wait, tick)
	})

	t.Run("should flush on max wait timeout", func(t *testing.T) {
		var counter int
		f := func() { counter++ }
		wait := 1 * time.Second
		tick := 100 * time.Millisecond
		options := NewDebounceOptions()
		options.MaxWait = 3 * wait
		debounced, cancel := Debounce(f, wait, options)
		defer cancel()

		require.Never(t, func() bool {
			debounced()
			return counter > 0
		}, 3*wait-tick, tick)
		require.Eventually(t, func() bool {
			return counter == 1
		}, 3*tick, tick)
	})
}
