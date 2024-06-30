package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
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

func promptConfirm(question, Y, N string) bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("%s (%s/%s): ", question, Y, N)
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input:", err)
		return false
	}

	input = strings.TrimSpace(input)
	return strings.EqualFold(input, Y)
}
