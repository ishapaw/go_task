package main

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/url"
	"time"
)

func main() {

	urls := []string{
		"https://google.com",
		"https://github.com",
		"https://stackoverflow.com",
		"https://www.wikipedia.org",
		"https://www.cloudflare.com",
		"http://example.com",
		"http://httpbin.org/status/200",
		"http://httpbin.org/status/404",
		"https://httpbin.org/delay/3",
		"https://www.reddit.com",
	}

	timeout := 5 * time.Second

	for _, rawURL := range urls {

		fmt.Print("Checking:", rawURL)

		u, err := url.Parse(rawURL)
		if err != nil {
			fmt.Println("Invalid URL:", err)
			continue
		}

		host := u.Hostname()
		scheme := u.Scheme

		if scheme == "http" {

			conn, err := net.DialTimeout("tcp", host+":80", timeout)
			if err != nil {
				fmt.Println("HTTP TCP connection failed:", err)
				return
			}

			fmt.Println("HTTP TCP connection successful")
			conn.Close()

		} else if scheme == "https" {

			conn, err := net.DialTimeout("tcp", host+":443", timeout)
			if err != nil {
				fmt.Println("HTTPS TCP connection failed:", err)
				return
			}

			tlsConn := tls.Client(conn, &tls.Config{
				ServerName: host,
			})

			err = tlsConn.Handshake()
			if err != nil {
				fmt.Println("TLS handshake failed:", err)
				conn.Close()
				return
			}

			fmt.Println("HTTPS TLS connection successful")
			tlsConn.Close()
		} else {
			fmt.Println("Unsupported scheme")
		}

	}

}
