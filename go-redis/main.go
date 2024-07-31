package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/deepanshuemblinux/go-redis/client"
	"github.com/deepanshuemblinux/go-redis/server"
	"github.com/sirupsen/logrus"
)

func main() {
	if len(os.Args) < 2 {
		logrus.Fatal("Usage: go-redis [server|client]")
	}
	switch strings.ToLower(os.Args[1]) {
	case "server":
		s, err := server.NewServer()
		if err != nil {
			logrus.WithField("Err", err).Fatal()
		}
		err = s.Run()
		if err != nil {
			logrus.WithField("Err", err).Fatal()
		}
	case "client":
		c, err := client.NewClient()
		if err != nil {
			logrus.WithField("Err", err).Fatal()
		}
		c.Init()
		var cmd string
		for {
			fmt.Scanln(&cmd)
			if strings.ToLower(cmd) == "exit" {
				os.Exit(0)
			}
			resp, err := c.Execute(cmd)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(resp)
		}
	}
}
