package network

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
)

func SendGET(conn net.Conn, host, path string) (int, error) {

	request := fmt.Sprintf(
		"GET %s HTTP/1.1\r\nHost: %s\r\nUser-Agent: go-health-checker\r\nConnection: close\r\n\r\n",
		path,
		host,
	)

	_, err := conn.Write([]byte(request))
	if err != nil {
		fmt.Println("Failed to send request:", err)
		return 0, err
	}

	reader := bufio.NewReader(conn)

	statusLine, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Failed to read response:", err)
		return 0, err
	}

	fmt.Println("Response:", statusLine)

	parts := strings.Split(statusLine, " ")
	if len(parts) < 2 {
		return 0, fmt.Errorf("invalid HTTP response")
	}

	code, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, err
	}

	return code, nil

}
