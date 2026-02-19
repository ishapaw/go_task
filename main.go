package main

import (
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

		fmt.Print("Checking")
	

	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		fmt.Println("parse error ", err)
	}

	host := parsedURL.Host

	fmt.Println(host)

	if parsedURL.Port() == "" {
		if parsedURL.Scheme == "https" {
			host += ":443"
		} else {
			host += ":80"
		}
	}

	conn, err := net.Dial("tcp", host)
	if err != nil {
		fmt.Println("connection error:", err)
		return
	}

	defer conn.Close()

	fmt.Println("Connected to", conn.RemoteAddr())
}
