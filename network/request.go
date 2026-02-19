package network

import (
	"bufio"
	"fmt"
	"net"
)

func SendGET(conn net.Conn, host, path string) {

	request := fmt.Sprintf(
		"GET %s HTTP/1.1\r\nHost: %s\r\nUser-Agent: go-health-checker\r\nConnection: close\r\n\r\n",
		path,
		host,
	)
	
	_, err := conn.Write([]byte(request))
	if err != nil {
		fmt.Println("Failed to send request:", err)
		return
	}

	reader := bufio.NewReader(conn)

	statusLine, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Failed to read response:", err)
		return
	}

	fmt.Println("Response:",statusLine)
}
