package cancellation

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSource(t *testing.T) {
	src := NewSource()

	assert.False(t, src.IsCancelled())

	select {
	case _, ok := <-src.Channel():
		assert.Fail(t, "Channel val, ok: %v", ok)
	default:
	}

	assert.True(t, src.Cancel())
	assert.True(t, src.IsCancelled())

	select {
	case _, ok := <-src.Channel():
		if ok {
			assert.Fail(t, "Channel has val")
		}
	default:
		assert.Fail(t, "Channel not closed")
	}

	assert.False(t, src.Cancel())
}

func TestSourceWithTimeout(t *testing.T) {
	src := NewSourceWithTimeout(500 * time.Millisecond)

	assert.False(t, src.IsCancelled())

	time.Sleep(time.Second)
	assert.True(t, src.IsCancelled())
}

func TestSourceParallel(t *testing.T) {
	var wg sync.WaitGroup

	ch := make(chan bool, 20)

	src := NewSource()
	assert.False(t, src.IsCancelled())

	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			assert.NotPanics(t, func() {
				ch <- src.Cancel()
			})
		}()
	}

	wg.Wait()
	close(ch)

	assert.True(t, src.IsCancelled())

	trueCount := 0
	for val := range ch {
		if val {
			trueCount++
		}
	}

	assert.Equal(t, 1, trueCount)
}

func TestCancellation(t *testing.T) {
	src := NewSource()
	c := src.Cancellation()

	assert.False(t, src.IsCancelled())
	assert.False(t, c.IsCancelled())

	select {
	case _, ok := <-c.Channel():
		assert.Fail(t, "Channel val, ok: %v", ok)
	default:
	}

	src.Cancel()
	assert.True(t, src.IsCancelled())
	assert.True(t, c.IsCancelled())

	select {
	case _, ok := <-c.Channel():
		if ok {
			assert.Fail(t, "Channel has val")
		}
	default:
		assert.Fail(t, "Channel not closed")
	}
}
