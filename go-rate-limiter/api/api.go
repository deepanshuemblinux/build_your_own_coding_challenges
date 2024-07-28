package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/deepanshuemblinux/go-rate-limiter/ratelimiter"
	"github.com/deepanshuemblinux/go-rate-limiter/service"
)

type apiServer struct {
	listenAddr       string
	srvc             service.MessageService
	rate_limiter_map map[string]ratelimiter.Ratelimiter
	limiter_type     int
}

func NewAPIServer(listenAddr string, srvc service.MessageService, limiter_type int) *apiServer {
	return &apiServer{
		listenAddr:       listenAddr,
		srvc:             srvc,
		rate_limiter_map: make(map[string]ratelimiter.Ratelimiter, 0),
		limiter_type:     limiter_type,
	}
}

func (s *apiServer) Run() {
	http.HandleFunc("/limited", s.handleLimited)
	http.HandleFunc("/unlimited", s.handleUnlimited)
	fmt.Printf("API Server listening on %s\n", s.listenAddr)
	err := http.ListenAndServe(s.listenAddr, nil)
	if err != nil {
		fmt.Println(err)
	}
}

func (s *apiServer) handleLimited(w http.ResponseWriter, r *http.Request) {
	ip := strings.Split(r.RemoteAddr, ":")[0]
	//log.Printf("Request came from %s\n", ip)
	log.Printf("Request came from %s\n", r.RemoteAddr)
	_, ok := s.rate_limiter_map[ip]
	if !ok {
		var limiter ratelimiter.Ratelimiter
		switch s.limiter_type {
		case ratelimiter.TokenBucket:
			limiter = ratelimiter.NewTokenBucket(10)
		case ratelimiter.FixedWindowCounter:
			limiter = ratelimiter.NewFixedWindowCounter(60, 60)
		case ratelimiter.SlidingWindowLog:
			limiter = ratelimiter.NewSlidingWindowLog(10, 30)
		}

		s.rate_limiter_map[ip] = limiter
		go s.rate_limiter_map[ip].StartLimiting()
	}
	if !s.rate_limiter_map[ip].IsAllowed() {
		fmt.Println("Rejecting request")
		w.WriteHeader(http.StatusTooManyRequests)
		return
	}
	resp := s.srvc.GetMessage("Limited, don't over use me!")
	err := writeJSON(w, http.StatusOK, resp)
	if err != nil {
		log.Println(err)
	}
}

func (s *apiServer) handleUnlimited(w http.ResponseWriter, r *http.Request) {
	resp := s.srvc.GetMessage("Unlimited, Let's go!")
	err := writeJSON(w, http.StatusOK, resp)
	if err != nil {
		log.Println(err)
	}
}

func writeJSON(w http.ResponseWriter, status int, value any) error {
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(value)
	if err != nil {
		return err
	}
	return nil
}
