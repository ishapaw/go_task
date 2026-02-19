package network

import (
	"crypto/tls"
	"net"
	"time"
)

func GetHTTPConnection(host string, timeout time.Duration) (net.Conn, error) {
	conn, err := net.DialTimeout("tcp", host+":80", timeout)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func GetHTTPSConnection(host string, timeout time.Duration) (net.Conn, error) {
	conn, err := net.DialTimeout("tcp", host+":443", timeout)
	if err != nil {
		return nil, err
	}

	tlsConn := tls.Client(conn, &tls.Config{
		ServerName: host,
	})

	err = tlsConn.Handshake()
	if err != nil {
		conn.Close()
		return nil, err
	}

	return tlsConn, nil
}
