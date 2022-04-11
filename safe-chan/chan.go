package safe_chan

import (
	"errors"
	"sync/atomic"
	"time"
)

var (
	ErrTimeoutExpired = errors.New("write timeout expired")
	ErrChanClosed     = errors.New("write to closed chan")
)

type SafeChan[E any] struct {
	C       chan E
	Options Options
	Closed  *int32
}

type Options struct {
	Retries int
	Timeout time.Duration
}

func NewSafeChan[E any](ch chan E, options Options) SafeChan[E] {
	return SafeChan[E]{
		C:       ch,
		Options: options,
		Closed:  new(int32),
	}
}

func (sc SafeChan[E]) writeWithTicker(d E, timeout time.Duration) error {
	t := time.NewTimer(timeout)
	defer t.Stop()

	select {
	case <-t.C:
	case sc.C <- d:
		return nil
	}

	return ErrTimeoutExpired
}

func (sc SafeChan[E]) write(d E, options Options) error {
	for i := 0; i < options.Retries; i++ {
		if err := sc.writeWithTicker(d, options.Timeout); err == nil {
			return nil
		}
	}

	return ErrTimeoutExpired
}

func (sc SafeChan[E]) Write(d E) error {
	if atomic.LoadInt32(sc.Closed) != 0 {
		return ErrChanClosed
	}
	return sc.write(d, sc.Options)
}

func (sc SafeChan[E]) Close() {
	atomic.StoreInt32(sc.Closed, 1)
	close(sc.C)
}

func (sc SafeChan[E]) WriteWithOptions(d E, options Options) error {
	return sc.write(d, options)
}
