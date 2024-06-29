package main

import (
	"fmt"
	"net"
	"time"
)

// pingHost tries to establish a TCP connection to the host to check accessibility
func pingHost(host string, port int, timeout time.Duration) error {
	hostaddr := fmt.Sprintf("%s:%d", host, port)
	conn, err := net.DialTimeout("tcp", hostaddr, timeout)
	if err != nil {
		return err
	}
	conn.Close()
	return nil
}
