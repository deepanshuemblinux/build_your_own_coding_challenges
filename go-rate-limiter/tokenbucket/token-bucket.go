package tokenbucket

import (
	"fmt"
	"sync"
	"time"
)

type TokenBucket interface {
	StartPushing()
	push(any)
	GetToken() bool
	StopPushing()
}

type defaultTokenBucket struct {
	len      int64
	mu       sync.Mutex
	stackBuf []any
	stopCh   chan bool
}

func NewTokenBucket(len int64) *defaultTokenBucket {
	t := defaultTokenBucket{len: len, mu: sync.Mutex{}}
	t.stackBuf = make([]any, 10)
	t.stopCh = make(chan bool)
	return &t
}

func (bucket *defaultTokenBucket) push(val any) {
	bucket.mu.Lock()
	defer bucket.mu.Unlock()
	if len(bucket.stackBuf) == int(bucket.len) {
		return
	}
	bucket.stackBuf = append(bucket.stackBuf, val)
}

func (bucket *defaultTokenBucket) GetToken() bool {
	bucket.mu.Lock()
	defer bucket.mu.Unlock()
	if len(bucket.stackBuf) == 0 {
		return false
	}
	fmt.Printf("Len of bucket is %d\n", len(bucket.stackBuf))
	bucket.stackBuf = bucket.stackBuf[:len(bucket.stackBuf)-1]
	return true
}

func (bucket *defaultTokenBucket) StartPushing() {
	for {
		select {
		case <-bucket.stopCh:
			return
		case <-time.NewTimer(time.Second).C:
			bucket.push(struct{}{})
		}
	}
}

func (bucket defaultTokenBucket) StopPushing() {
	bucket.stopCh <- true
}
