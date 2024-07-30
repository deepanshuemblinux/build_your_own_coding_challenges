package loadbalancer

import (
	"io"
	"log"
	"net/http"
	"slices"
	"sync"
	"time"

	"github.com/deepanshuemblinux/build_your_own_coding_challenges/go-load-balancer/httpserver"
	"github.com/sirupsen/logrus"
)

type LoadBalancingAlgo interface {
	GetNextBackend(lb *Lb) (string, error)
}

type Lb struct {
	httpserver.Server
	healthy_backends []string
	all_backends     []string
	client           *http.Client
	algo             LoadBalancingAlgo
	Mu_Healthy       sync.Mutex
	Mu               sync.Mutex
}

func (l *Lb) GetHealthyBackends() []string {
	var healthy_backends_copy []string
	l.Mu_Healthy.Lock()
	defer l.Mu_Healthy.Unlock()
	healthy_backends_copy = append(healthy_backends_copy, l.healthy_backends...)
	return healthy_backends_copy
}

func (l *Lb) GetAllBackends() []string {
	var backends_copy []string
	l.Mu.Lock()
	defer l.Mu.Unlock()
	backends_copy = append(backends_copy, l.all_backends...)
	return backends_copy
}
func (l *Lb) requestHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logrus.WithFields(logrus.Fields{
			"From":       r.RemoteAddr,
			"User Agent": r.UserAgent(),
			"Method":     r.Method,
			"Protocol":   r.Proto,
		}).Info("Recieved Request")

		backend_url, err := l.algo.GetNextBackend(l)
		logrus.WithField("Forwarding request to:", backend_url)
		if err != nil {
			logrus.WithField("Backend:", backend_url).Error(err)
			w.WriteHeader(http.StatusBadGateway)
			return
		}
		breq, err := http.NewRequest(r.Method, backend_url, r.Body)
		if err != nil {
			log.Printf("Error creating request object %s\n", err)
		}
		resp, err := l.client.Do(breq)
		if err != nil {
			log.Printf("Error getting response from %s\n", backend_url)
		}
		w.WriteHeader(http.StatusOK)
		respBody, _ := io.ReadAll(resp.Body)
		w.Write(respBody)
	}
}
func NewLB(listenAddr string, algo LoadBalancingAlgo) *Lb {
	lb := Lb{
		httpserver.Server{
			ListenAddr: listenAddr,
			Mux:        http.NewServeMux(),
		},
		make([]string, 0),
		make([]string, 0),
		http.DefaultClient,
		algo,
		sync.Mutex{},
		sync.Mutex{},
	}

	lb.Mux.HandleFunc("/", lb.requestHandler())
	go lb.healthCheck()
	return &lb
}

func (l *Lb) RegisterBackend(backend string) {
	l.Mu.Lock()
	defer l.Mu.Unlock()
	l.all_backends = append(l.all_backends, backend)
}

func (l *Lb) healthCheck() {
	for {
		time.Sleep(time.Second * 10)
		backends := l.GetAllBackends()
		for _, backend_url := range backends {
			breq, err := http.NewRequest("GET", backend_url, nil)
			if err != nil {
				log.Printf("Error creating request object %s\n", err)
			}
			resp, err := l.client.Do(breq)
			l.Mu_Healthy.Lock()
			idx := slices.Index(l.healthy_backends, backend_url)
			if err != nil && idx != -1 {
				log.Printf("Removing unhealthy backend %s\n", backend_url)
				l.healthy_backends = append(l.healthy_backends[0:idx], l.healthy_backends[idx+1:]...)
			} else if err == nil && resp.StatusCode == http.StatusOK && idx == -1 {
				log.Printf("Adding back healthy backend %s\n", backend_url)
				l.healthy_backends = append(l.healthy_backends, backend_url)
			} else if err == nil && resp.StatusCode != http.StatusOK && idx != -1 {
				log.Printf("Removing unhealthy backend %s\n", backend_url)
				l.healthy_backends = append(l.healthy_backends[0:idx], l.healthy_backends[idx+1:]...)
			}
			l.Mu_Healthy.Unlock()
		}

	}
}
