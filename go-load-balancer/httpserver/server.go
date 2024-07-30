package httpserver

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

type Server struct {
	ListenAddr string
	Mux        *http.ServeMux
	Name       string
}

func NewServer(listenAddr, name string) *Server {
	return &Server{
		ListenAddr: listenAddr,
		Mux:        http.NewServeMux(),
		Name:       name,
	}
}

func (s *Server) GetListenAddr() string {
	return s.ListenAddr
}
func (s *Server) Run() {
	logrus.WithField("Address", s.ListenAddr).Info("Server started listeining")
	http.ListenAndServe(s.ListenAddr, s.Mux)
}

func (s *Server) HandleFunc(path string, f http.HandlerFunc) {

	s.Mux.HandleFunc(path, f)
}
