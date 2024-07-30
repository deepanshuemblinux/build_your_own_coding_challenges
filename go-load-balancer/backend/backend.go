package backend

import (
	"net/http"

	"github.com/deepanshuemblinux/build_your_own_coding_challenges/go-load-balancer/httpserver"
)

type Backend struct {
	httpserver.Server
}

func NewBackend(listenAddr string) *Backend {
	return &Backend{
		httpserver.Server{
			ListenAddr: listenAddr,
			Mux:        http.NewServeMux(),
		},
	}
}
