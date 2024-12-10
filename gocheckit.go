package gocheckit

import (
	"net/http"
	"sync"
)

// Client allow caller override the default client
var Client = http.DefaultClient

// Do send HEAD requests to domain in slice, n domain at the same time.
func Do(domains []string, n int) bool {
	limitChan := make(chan struct{}, n)
	resultChan := make(chan bool)
	var wg sync.WaitGroup
	checkit := func(domain string) {
		defer wg.Done()

		limitChan <- struct{}{}
		defer func() { <-limitChan }() // unblock the channel

		res, err := Client.Head(domain)
		if err == nil && res.StatusCode == http.StatusOK {
			resultChan <- true
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

	for res := range resultChan {
		if res {
			return true
		}
	}

	return false
}
