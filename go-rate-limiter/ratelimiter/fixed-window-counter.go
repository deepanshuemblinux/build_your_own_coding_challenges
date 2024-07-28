package ratelimiter

import (
	"sync"
	"time"
)

type fixedWindowCounter struct {
	mu         sync.Mutex
	limit      int
	windowSize int
	stopCh     chan bool
	counter    int
	enabled    bool
}

func NewFixedWindowCounter(limit, windowSize int) *fixedWindowCounter {
	return &fixedWindowCounter{
		limit:      limit,
		windowSize: windowSize,
		mu:         sync.Mutex{},
		stopCh:     make(chan bool),
	}
}

func (fwc *fixedWindowCounter) IsAllowed() bool {
	if !fwc.enabled {
		return true
	}
	fwc.mu.Lock()
	defer fwc.mu.Unlock()
	if fwc.counter >= fwc.limit {
		return false
	}
	fwc.counter++
	return true
}

func (fwc *fixedWindowCounter) StartLimiting() {
	fwc.enabled = true
	for {
		select {
		case <-fwc.stopCh:
			fwc.enabled = false
			return
		case <-time.NewTimer(time.Second * time.Duration(fwc.windowSize)).C:
			fwc.mu.Lock()
			fwc.counter = 0
			fwc.mu.Unlock()
		}
	}
}

func (fwc *fixedWindowCounter) StopLimiting() {
	fwc.stopCh <- true

}
