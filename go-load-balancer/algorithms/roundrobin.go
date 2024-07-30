package algorithms

import (
	"fmt"

	"github.com/deepanshuemblinux/build_your_own_coding_challenges/go-load-balancer/loadbalancer"
)

type RoundRobin struct {
	next_index int
}

func NewRoundRobin() *RoundRobin {
	return &RoundRobin{
		next_index: -1,
	}
}
func (rr *RoundRobin) GetNextBackend(lb *loadbalancer.Lb) (string, error) {
	backends := lb.GetHealthyBackends()
	fmt.Println("Healthy backend is get next backend", backends)
	if len(backends) == 0 {
		return "", fmt.Errorf("No healthy backends Found")
	}
	rr.next_index = (rr.next_index + 1) % len(backends)
	fmt.Println("next available index is ", rr.next_index)
	fmt.Println("next available backend is ", backends[rr.next_index])
	return backends[rr.next_index], nil
}
