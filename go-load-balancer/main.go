package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/deepanshuemblinux/build_your_own_coding_challenges/go-load-balancer/algorithms"
	"github.com/deepanshuemblinux/build_your_own_coding_challenges/go-load-balancer/backend"
	"github.com/deepanshuemblinux/build_your_own_coding_challenges/go-load-balancer/loadbalancer"
	"github.com/sirupsen/logrus"
)

func handleBackend(w http.ResponseWriter, r *http.Request) {
	logrus.WithFields(logrus.Fields{
		"From":       r.RemoteAddr,
		"User Agent": r.UserAgent(),
		"Method":     r.Method,
		"Protocol":   r.Proto,
	}).Info("Recieved Request")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Hello from backend %s", r.Host)))
	log.Println("Replied with a hello message")
}

func main() {
	if len(os.Args) < 3 {
		log.Fatal("Usage: go run main.go [backend|lb] [listen address]")
	}
	ch := make(chan bool)
	switch strings.ToUpper(os.Args[1]) {
	case "BACKEND":

		backend := backend.NewBackend(os.Args[2])
		backend.HandleFunc("/", handleBackend)
		backend.Run()

		<-ch
	case "LB":
		lb := loadbalancer.NewLB(os.Args[2], algorithms.NewRoundRobin())
		lb.RegisterBackend("http://localhost:8081/")
		lb.RegisterBackend("http://localhost:8082/")
		lb.RegisterBackend("http://localhost:8083/")
		lb.Run()
		<-ch
	}
}
