package matcha

import (
	"time"
)

type throttler struct {
	ticker *time.Ticker
	ch     chan struct{}
}

func newThrottler(fps int) *throttler {
	throttler := &throttler{
		ticker: time.NewTicker(time.Second / time.Duration(fps)),
		ch:     make(chan struct{}, 1),
	}

	go func(ticker *time.Ticker, ch chan<- struct{}) {
		for range ticker.C {
			select {
			case ch <- struct{}{}:
			default:
			}
		}
	}(throttler.ticker, throttler.ch)
	return throttler
}

func (t *throttler) trigger() {
	select {
	case t.ch <- struct{}{}:
	default:
	}
}

func (t *throttler) channel() <-chan struct{} {
	return t.ch
}
