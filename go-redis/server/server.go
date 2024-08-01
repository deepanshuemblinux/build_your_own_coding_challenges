package server

import (
	"encoding/json"
	"fmt"
	"net"
	"os"

	"github.com/deepanshuemblinux/go-redis/utils"
	"github.com/sirupsen/logrus"
)

type Server struct {
	listenAddr string
	stopCh     chan bool
}

func NewServer() (*Server, error) {
	config_file, err := os.Open("config/server.json")
	if err != nil {
		return nil, err
	}
	config := map[string]interface{}{}
	json.NewDecoder(config_file).Decode(&config)
	address, ok := config["address"]
	if !ok {
		return nil, fmt.Errorf("provide value for adrress in config.json")
	}

	port, ok := config["port"]
	if !ok {
		return nil, fmt.Errorf("provide value for port in config.json")
	}

	listenAddr := fmt.Sprintf("%s:%s", address, port)
	return &Server{
		listenAddr: listenAddr,
		stopCh:     make(chan bool),
	}, nil
}

func (s *Server) Run() error {
	listener, err := net.Listen("tcp", s.listenAddr)
	if err != nil {
		return err
	}
	logrus.WithFields(logrus.Fields{
		"Address": listener.Addr().String(),
	}).Info("Running Redis Server")
	defer listener.Close()
	go s.acceptLoop(listener)
	<-s.stopCh
	return nil
}

func (s *Server) acceptLoop(listener net.Listener) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			logrus.WithField("err:", err).Fatal("Accepting Connection Failed")
			continue
		}
		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	req := make([][]byte, 0)
	utils.ParseRequest(conn, &req)
	for i, r := range req {
		fmt.Printf("%dth index of request array is %s\n", i, string(r))
	}
	resp, err := utils.HandleCommand(&req)
	if err != nil {
		logrus.Error(err)
	}
	n1, err := conn.Write(resp)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("num bytes written is ", n1)

}

func (s *Server) Stop() {
	s.stopCh <- true
}
