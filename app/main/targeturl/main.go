package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"
)

func main() {
	// Validate input: a URL must be supplied as a command-line argument.
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <url>")
		return
	}
	originalURL := os.Args[1]

	// Instantiate an HTTP client with a custom redirection policy.
	// The CheckRedirect function prevents automatic following of redirects,
	// thereby allowing us to capture and examine the "Location" header.
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	// Execute an HTTP GET request on the provided URL.
	response, err := client.Get(originalURL)
	if err != nil {
		fmt.Printf("Error retrieving URL: %v\n", err)
		return
	}
	defer response.Body.Close()

	// Check if the HTTP response indicates a redirection (3xx status code).
	if response.StatusCode >= http.StatusMultipleChoices && response.StatusCode < http.StatusBadRequest {
		// Extract the target URL from the "Location" header.
		targetURL := response.Header.Get("Location")
		// Parse the target URL to enable structured manipulation.
		parsedURL, err := url.Parse(targetURL)
		if err != nil {
			fmt.Printf("Error parsing target URL: %v\n", err)
			return
		}

		// Remove the "ref=producthunt" parameter from the query if present.
		queryParams := parsedURL.Query()
		if queryParams.Get("ref") == "producthunt" {
			queryParams.Del("ref")
			parsedURL.RawQuery = queryParams.Encode()
		}
		finalURL := parsedURL.String()
		fmt.Printf("The provided URL redirects to: %s\n", finalURL)

		// Check if the target URL is live and print the result.
		liveStatus := checkURLLive(finalURL)
		fmt.Println(liveStatus)
	} else {
		// Inform the user if no redirection is detected.
		fmt.Printf("The provided URL does not result in a redirection (Status code: %d).\n", response.StatusCode)
	}
}

// checkURLLive performs a HEAD request to determine if the target URL is accessible.
func checkURLLive(target string) string {
	// Create an HTTP client with a timeout to avoid prolonged waits.
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	// Perform a HEAD request to minimize data transfer.
	resp, err := client.Head(target)
	if err != nil {
		return fmt.Sprintf("Website is down: %v", err)
	}
	defer resp.Body.Close()

	// Determine the liveliness based on the HTTP status code.
	if resp.StatusCode >= 200 && resp.StatusCode < 400 {
		return "Website is live."
	}
	return fmt.Sprintf("Webpage is unavailable (Status code: %d).", resp.StatusCode)
}
