package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"

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
	req := bytes.NewBuffer([]byte{})
	n, err := io.Copy(req, conn)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("num bytes read is ", n)
	fmt.Println("request from client is ", string(req.Bytes()))
	n1, err := conn.Write([]byte("+PONG\r\n"))
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("num bytes written is ", n1)

}

func (s *Server) Stop() {
	s.stopCh <- true
}
