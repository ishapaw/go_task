package main

import (
	"fmt"
	"net/url"
	"time"

	"go_task/network"
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
		"https://httpbin.org/delay/3",
		"https://www.reddit.com",
	}

	timeout := 5 * time.Second

	for _, rawURL := range urls {
		fmt.Println("Checking:", rawURL)

		u, err := url.Parse(rawURL)
		if err != nil {
			fmt.Println("Invalid URL:", err)
			continue
		}

		host := u.Hostname()
		path := u.Path
		if path == "" {
			path = "/"
		}

		var connErr error

		if u.Scheme == "http" {
			conn, err := network.GetHTTPConnection(host, timeout)

			if err == nil {
				network.SendGET(conn, host, path)
				conn.Close()
			}

			connErr = err

		} else if u.Scheme == "https" {

			conn, err := network.GetHTTPSConnection(host, timeout)

			if err == nil {
				network.SendGET(conn, host, path)
				conn.Close()
			}

			connErr = err
		}

		if connErr != nil {
			fmt.Println("Connection failed:", connErr)
		}

		fmt.Println()
	}
}
