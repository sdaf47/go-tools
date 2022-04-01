package safe_chan

import (
	"errors"
	"testing"
	"time"
)

func TestSafeChan_Write(t *testing.T) {
	const (
		retries    = 1
		timeout    = 100 * time.Millisecond
		bufferSize = 2
	)

	ch := NewSafeChan(make(chan int, bufferSize), Options{
		Retries: retries,
		Timeout: timeout,
	})

	err := ch.Write(1)
	if err != nil {
		t.Fatal(err)
	}

	err = ch.Write(1)
	if err != nil {
		t.Fatal(err)
	}

	err = ch.Write(2)
	if !errors.Is(err, ErrTimeoutExpired) {
		t.Fatalf("expect: %s, got: %s", ErrTimeoutExpired, err)
	}
}

func TestSafeChan_Close(t *testing.T) {
	const (
		retries    = 1
		timeout    = 100 * time.Millisecond
		bufferSize = 2
	)

	ch := NewSafeChan(make(chan int, bufferSize), Options{
		Retries: retries,
		Timeout: timeout,
	})

	err := ch.Write(1)
	if err != nil {
		t.Fatal(err)
	}

	ch.Close()

	err = ch.Write(2)
	if !errors.Is(err, ErrChanClosed) {
		t.Fatalf("expect: %s, got: %s", ErrTimeoutExpired, err)
	}
}
