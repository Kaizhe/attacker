package utils

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"strconv"
)

// LogPrint print an response
func LogPrint(w http.ResponseWriter, msg string, status int) {
	fmt.Println(msg)
	w.WriteHeader(status)
	w.Write([]byte(msg + "\n"))
}

// GetFreePort returns next available port to reverse shell contact
func GetFreePort() (string, error) {
	addr, err := net.ResolveTCPAddr("tcp", "0.0.0.0:0")
	if err != nil {
		return "", errors.New("fail to get a free port")
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return "", errors.New("fail to test on free port")
	}
	defer l.Close()
	return strconv.Itoa(l.Addr().(*net.TCPAddr).Port), nil
}
