package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"

	"github.com/sirupsen/logrus"
)

type Client struct {
	serveraddr string
	conn       net.Conn
}

func NewClient() (*Client, error) {
	config_file, err := os.Open("config/client.json")
	if err != nil {
		return nil, fmt.Errorf("error opening config file %w", err)
	}
	config := map[string]interface{}{}
	err = json.NewDecoder(config_file).Decode(&config)
	if err != nil {
		return nil, fmt.Errorf("error decoding config json %w", err)
	}
	address, ok := config["address"]
	if !ok {
		return nil, fmt.Errorf("provide value for adrress in config.json")
	}

	port, ok := config["port"]
	if !ok {
		return nil, fmt.Errorf("provide value for port in config.json")
	}

	serverAddr := fmt.Sprintf("%s:%s", address, port)
	return &Client{
		serveraddr: serverAddr,
	}, nil
}

func (c *Client) Init() {
	conn, err := net.Dial("tcp", c.serveraddr)
	if err != nil {
		logrus.WithField("Server Adress", c.serveraddr).Fatal("Error connecting to Server")
	}
	c.conn = conn
}
func (c *Client) Execute(cmd string) (string, error) {
	cmd_bytes := []byte(cmd)
	_, err := c.conn.Write(cmd_bytes)
	if err != nil {
		return "", fmt.Errorf("error executing redis command %w", err)
	}
	resp_buff := bytes.NewBuffer([]byte{})
	_, err = io.Copy(resp_buff, c.conn)
	if err != nil {
		return "", fmt.Errorf("error getting response for redis command %w", err)
	}
	return string(resp_buff.Bytes()), nil
}
