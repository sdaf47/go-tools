package safe_chan

import (
	"errors"
	"sync/atomic"
	"time"
)

var (
	ErrTimeoutExpired = errors.New("channel write timeout expired")
	ErrChanClosed     = errors.New("write to closed chan")
)

type SafeChan[E any] struct {
	channel chan E
	options Options
	closed  *int32
}

type Options struct {
	retries int
	timeout time.Duration
}

func NewSafeChan[E any](ch chan E, options Options) SafeChan[E] {
	return SafeChan[E]{
		channel: ch,
		options: options,
		closed:  new(int32),
	}
}

func (sc SafeChan[E]) write(d E, options Options) error {
	t := time.NewTimer(options.timeout)
	defer t.Stop()

	for i := 0; i < options.retries; i++ {
		select {
		case <-t.C:
			continue
		case sc.channel <- d:
			return nil
		}
	}

	return ErrTimeoutExpired
}

func (sc SafeChan[E]) Write(d E) error {
	if atomic.LoadInt32(sc.closed) != 0 {
		return ErrChanClosed
	}
	return sc.write(d, sc.options)
}

func (sc SafeChan[E]) C() chan E {
	return sc.channel
}

func (sc SafeChan[E]) Close() {
	atomic.StoreInt32(sc.closed, 1)
	close(sc.channel)
}

func (sc SafeChan[E]) WriteWithOptions(d E, options Options) error {
	return sc.write(d, options)
}
