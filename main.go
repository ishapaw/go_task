package main

import (
	"fmt"
	"net/url"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"go_task/network"
)

type Result struct {
	URL        string
	Status     string
	StatusCode int
	Error      error
}

func checkURL(rawURL string, timeout time.Duration) Result {
	u, err := url.Parse(rawURL)
	if err != nil {
		return Result{
			URL:    rawURL,
			Status: "DOWN",
			Error:  fmt.Errorf("invalid URL: %w", err),
		}
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return Result{
			URL:    rawURL,
			Status: "DOWN",
			Error:  fmt.Errorf("unsupported scheme: %s", u.Scheme),
		}
	}

	host := u.Hostname()
	port := u.Port()

	path := u.RequestURI()
	if path == "" {
		path = "/"
	}

	var connErr error
	var statusCode int

	if u.Scheme == "http" {
		conn, err := network.GetHTTPConnection(host, port, timeout)
		if err == nil {
			statusCode, err = network.SendGET(conn, host, path)
			conn.Close()
		}
		connErr = err

	} else {
		conn, err := network.GetHTTPSConnection(host, port, timeout)
		if err == nil {
			statusCode, err = network.SendGET(conn, host, path)
			conn.Close()
		}
		connErr = err
	}

	if connErr != nil {
		return Result{
			URL:    rawURL,
			Status: "DOWN",
			Error:  connErr,
		}
	}

	if statusCode >= 500 {
		return Result{
			URL:        rawURL,
			Status:     "DOWN",
			StatusCode: statusCode,
		}
	}

	return Result{
		URL:        rawURL,
		Status:     "UP",
		StatusCode: statusCode,
	}
}

func runHealthChecks(urls []string, timeout time.Duration) []Result {
	results := make([]Result, len(urls))
	var wg sync.WaitGroup

	for i, rawURL := range urls {
		wg.Add(1)

		go func(idx int, u string) {
			defer wg.Done()
			results[idx] = checkURL(u, timeout)
		}(i, rawURL)
	}

	wg.Wait()
	return results
}

func printSummary(results []Result) {
	fmt.Println("\n========== Health Check Summary ==========")

	var upList, downList []Result

	for _, r := range results {
		if r.Status == "UP" {
			upList = append(upList, r)
		} else {
			downList = append(downList, r)
		}
	}

	fmt.Printf("\nUP (%d):\n", len(upList))
	for _, r := range upList {
		fmt.Printf("  ✓ %-45s [status %d]\n", r.URL, r.StatusCode)
	}

	fmt.Printf("\nDOWN (%d):\n", len(downList))
	for _, r := range downList {
		if r.Error != nil {
			fmt.Printf("  ✗ %-45s [error: %v]\n", r.URL, r.Error)
		} else {
			fmt.Printf("  ✗ %-45s [status %d]\n", r.URL, r.StatusCode)
		}
	}

	fmt.Println("\n==========================================")
}

func main() {
	urls := []string{
		"http://10.255.255.1",
		"https://github.com",
		"https://stackoverflow.com",
		"https://www.wikipedia.org",
		"https://www.cloudflare.com",
		"http://example.com",
		"http://httpbin.org/status/200",
		"https://httpbin.org/delay/3",
		"https://www.reddit.com",
		"ftp://example.com",
	}

	timeout := 5 * time.Second
	interval := 10 * time.Second

	fmt.Println("Starting HTTP Health Checker")
	fmt.Printf("Checking %d URLs every %v (timeout: %v)\n",
		len(urls), interval, timeout)
	fmt.Println("Press Ctrl+C to stop.")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	results := runHealthChecks(urls, timeout)
	printSummary(results)

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			fmt.Printf("\n[%s] Running health checks...\n",
				time.Now().Format("15:04:05"))

			results = runHealthChecks(urls, timeout)
			printSummary(results)

		case sig := <-quit:
			fmt.Printf("\n\nReceived %v signal. Shutting down gracefully...\n", sig)
			fmt.Println("Health checker stopped.")
			os.Exit(0)
		}
	}
} 