package utils

import (
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/deepanshuemblinux/go-redis/db"
	"github.com/sirupsen/logrus"
)

func discardCRLF(conn net.Conn) {
	buf := make([]byte, 2)
	conn.Read(buf)
}
func getAggregateSize(conn net.Conn, req *[][]byte) (size int, err error) {
	buf := make([]byte, 1)
	size_buf := make([]byte, 0)

	for {
		_, err = conn.Read(buf)
		if err != nil {
			return size, err
		}
		if buf[0] != '\r' {
			size_buf = append(size_buf, buf...)
		} else {
			buf := make([]byte, 1)
			conn.Read(buf)
			break
		}
	}
	size, err = strconv.Atoi(string(size_buf))
	if err != nil {
		return size, err
	}
	return size, nil
}
func handleArrays(conn net.Conn, req *[][]byte) error {
	size, err := getAggregateSize(conn, req)
	if err != nil {
		logrus.Error(err)
		return err
	}
	for i := 0; i < size; i++ {
		ParseRequest(conn, req)
	}
	return nil
}

func handleBulkStrings(conn net.Conn, req *[][]byte) error {
	size, err := getAggregateSize(conn, req)
	if err != nil {
		return err
	}
	buf := make([]byte, size)
	_, err = conn.Read(buf)
	if err != nil {
		return err
	}
	*req = append(*req, buf)
	discardCRLF(conn)
	return nil
}

func handleSet(req *[][]byte) []byte {
	db.KV_DB.Set(string((*req)[1]), string((*req)[2]))
	val, _ := db.KV_DB.Get(string((*req)[1]))
	logrus.Info("After setting the key value, the value is ", val)
	return []byte("+OK\r\n")
}
func handleBulkStrResp(req *[][]byte) []byte {
	respStr := "$"
	respStr = fmt.Sprintf("%s%d\r\n", respStr, len((*req)[1]))
	respStr = fmt.Sprintf("%s%s\r\n", respStr, (*req)[1])
	return []byte(respStr)
}
func HandleCommand(req *[][]byte) ([]byte, error) {
	switch strings.ToLower(string((*req)[0])) {
	case "ping":
		return []byte("+PONG\r\n"), nil
	case "echo":
		return handleBulkStrResp(req), nil
	case "set":
		return handleSet(req), nil
	}
	return nil, nil
}
func ParseRequest(conn net.Conn, req *[][]byte) error {
	data_type := make([]byte, 1)
	_, err := conn.Read(data_type)
	if err != nil {
		return fmt.Errorf("error parsing request %w", err)
	}
	switch string(data_type[0]) {
	case "+":

	case "-":

	case "$":
		handleBulkStrings(conn, req)
	case ":":

	case "*":
		handleArrays(conn, req)
	default:
		return fmt.Errorf("invalid data type")
	}
	return nil
}
