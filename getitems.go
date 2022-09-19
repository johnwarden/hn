package hn

import (
	"fmt"
	"sync"
	"sync/atomic"
)

const maxGoroutines = 100

// GetItems gets a set of from the HN API in parallel.
func (client *Client) GetItems(items []int) ([]Item, error) {

	n := len(items)
	results := make([]Item, n)

	// Use a channel as a rate limiter. If there are more than
	// nGoroutines running reading from the channel will block
	// until one goroutine releases.
	sem := make(chan struct{}, maxGoroutines)
	acquire := func() { sem <- struct{}{} }
	release := func() { <-sem }

	var wg sync.WaitGroup

	var nSuccess int64 = 0

	for index, itemID := range items {
		acquire()
		wg.Add(1)
		go func(idx, id int) {
			defer release()
			defer wg.Done()
			item, err := client.Item(id)
			if err != nil {
				fmt.Println("failed to fetch story", id, err)
				return
			}
			atomic.AddInt64(&nSuccess, 1)
			results[idx] = *item
		}(index, itemID)
	}

	wg.Wait()

	fmt.Println("Number of successes", nSuccess, "out of", n)
	if nSuccess != int64(n) {
		return results, fmt.Errorf("Didn't successfully fetch all items.")
	}

	return results, nil
}
