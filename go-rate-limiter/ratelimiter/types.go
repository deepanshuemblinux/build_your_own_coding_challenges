package ratelimiter

const (
	TokenBucket = iota
	FixedWindowCounter
	SlidingWindowLog
)

type Ratelimiter interface {
	StartLimiting()
	IsAllowed() bool
	StopLimiting()
}
