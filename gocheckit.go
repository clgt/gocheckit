package gocheckit

import (
	"net/http"
	"sync"
)

// Client allow caller override the default client
var Client = http.DefaultClient

// Do send HEAD requests to domain in slice, n domain at the same time.
func Do(domains []string, n int) (string, bool) {
	limitChan := make(chan struct{}, n)
	resultChan := make(chan string)
	var wg sync.WaitGroup
	checkit := func(domain string) {
		defer wg.Done()

		limitChan <- struct{}{}
		defer func() { <-limitChan }() // unblock the channel

		res, err := Client.Head(domain)
		if err == nil && res.StatusCode == http.StatusOK {
			resultChan <- domain
		}
	}

	for _, domain := range domains {
		wg.Add(1)
		go checkit(domain)
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	for domain := range resultChan {
		return domain, true
	}

	return "", false
}
