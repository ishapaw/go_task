package network

import (
	"crypto/tls"
	"fmt"
	"net"
	"time"
)

func GetHTTPConnection(host, port string, timeout time.Duration) (net.Conn, error) {
	if port == "" {
		port = "80"
	}

	conn, err := net.DialTimeout("tcp", host+":"+port, timeout)
	if err != nil {
		return nil, err
	}

	if err := conn.SetDeadline(time.Now().Add(timeout)); err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to set deadline: %w", err)
	}

	return conn, nil
}

func GetHTTPSConnection(host, port string, timeout time.Duration) (net.Conn, error) {
	if port == "" {
		port = "443"
	}

	conn, err := net.DialTimeout("tcp", host+":"+port, timeout)
	if err != nil {
		return nil, err
	}

	if err := conn.SetDeadline(time.Now().Add(timeout)); err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to set initial deadline: %w", err)
	}

	tlsConn := tls.Client(conn, &tls.Config{ServerName: host})

	err = tlsConn.Handshake()
	if err != nil {
		tlsConn.Close()
		return nil, err
	}

	if err := tlsConn.SetDeadline(time.Now().Add(timeout)); err != nil {
		tlsConn.Close()
		return nil, fmt.Errorf("failed to set post-handshake deadline: %w", err)
	}

	return tlsConn, nil
}
