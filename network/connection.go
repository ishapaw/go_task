package network

import (
	"crypto/tls"
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

	conn.SetDeadline(time.Now().Add(timeout))

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

	conn.SetDeadline(time.Now().Add(timeout))

	tlsConn := tls.Client(conn, &tls.Config{ServerName: host})

	err = tlsConn.Handshake()
	if err != nil {
		conn.Close()
		return nil, err
	}

	tlsConn.SetDeadline(time.Now().Add(timeout))
	
	return tlsConn, nil
}
