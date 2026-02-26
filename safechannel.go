package dag_go

import (
	"fmt"
	"sync"
)

// SafeChannel is a generic, concurrency-safe channel wrapper that prevents
// double-close panics and provides non-blocking send semantics.
type SafeChannel[T any] struct {
	ch     chan T
	closed bool
	mu     sync.RWMutex
}

// NewSafeChannelGen creates a new SafeChannel with the given buffer size.
func NewSafeChannelGen[T any](buffer int) *SafeChannel[T] {
	return &SafeChannel[T]{
		ch:     make(chan T, buffer),
		closed: false,
	}
}

// Send attempts to deliver value to the channel. Returns false if the channel
// is already closed or the buffer is full.
func (sc *SafeChannel[T]) Send(value T) bool {
	sc.mu.RLock()
	defer sc.mu.RUnlock()

	if sc.closed {
		return false
	}
	select {
	case sc.ch <- value:
		return true
	default:
		return false
	}
}

// Close closes the underlying channel exactly once. Returns an error if the
// channel is already closed or a panic occurs during close.
func (sc *SafeChannel[T]) Close() (err error) {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	if sc.closed {
		return fmt.Errorf("channel already closed")
	}

	// panic 이 발생하면 err 에 메시지를 저장한다.
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic while closing channel: %v", r)
		}
	}()

	close(sc.ch)
	sc.closed = true
	return nil
}

// GetChannel returns the underlying channel for range/select operations.
func (sc *SafeChannel[T]) GetChannel() chan T {
	return sc.ch
}
