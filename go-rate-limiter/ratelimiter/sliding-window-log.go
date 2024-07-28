package ratelimiter

import (
	"fmt"
	"sync"
	"time"
)

type slidingWindowLog struct {
	logsTable  []time.Time
	enabled    bool
	limit      int
	windowSize int
	mu         sync.Mutex
}

func NewSlidingWindowLog(limit, windowSize int) *slidingWindowLog {
	return &slidingWindowLog{
		logsTable:  make([]time.Time, 0),
		limit:      limit,
		windowSize: windowSize,
	}
}

func (swl *slidingWindowLog) StartLimiting() {
	swl.enabled = true
}

func (swl *slidingWindowLog) StopLimiting() {
	swl.enabled = false
}

func (swl *slidingWindowLog) IsAllowed() bool {
	if !swl.enabled {
		return true
	}
	currentTime := time.Now()

	//swl.mu.Lock()
	//defer swl.mu.Unlock()
	fmt.Printf("Length of log table is %d\n", len(swl.logsTable))
	fmt.Printf("Limit is %d\n", swl.limit)
	if len(swl.logsTable) >= swl.limit {
		fmt.Println("inside if")
		timeDiff := currentTime.Unix() - swl.logsTable[len(swl.logsTable)-swl.limit].Unix()
		fmt.Printf("Current time is %s\n", currentTime)
		fmt.Printf("timeDiff is %d\n", timeDiff)
		if timeDiff >= int64(swl.windowSize) {
			swl.logsTable = append(swl.logsTable, currentTime)
			swl.logsTable = swl.logsTable[1:]
			return true
		}
		return false
	} else {
		fmt.Println("inside else")
		swl.logsTable = append(swl.logsTable, currentTime)
		return true
	}
}
