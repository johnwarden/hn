package hn

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

const defaultMaxGoroutines = 20

// GetItems gets a set of from the HN API in parallel.
func (client *Client) GetItems(ctx context.Context, items []int, maxGoroutines int) ([]Item, error) {
	if maxGoroutines == 0 {
		maxGoroutines = defaultMaxGoroutines
	}

	n := len(items)
	results := make([]Item, n)

	// Use a channel as a rate limiter. If there are more than
	// nGoroutines running reading from the channel will block
	// until one goroutine releases.
	sem := make(chan struct{}, maxGoroutines)
	acquire := func() { sem <- struct{}{} }
	release := func() { <-sem }

	var wg sync.WaitGroup

	var nSuccess int64

	contextCanceled := false
STORIES:
	for index, itemID := range items {
		acquire()
		wg.Add(1)
		go func(idx, id int) {
			defer release()
			defer wg.Done()

			c, cancel := context.WithTimeout(ctx, 10*time.Second)
			defer cancel()

			item, err := client.Item(c, id)
			if err != nil {
				return
			}
			atomic.AddInt64(&nSuccess, 1)
			results[idx] = *item
		}(index, itemID)

		select {
		case <-ctx.Done():
			contextCanceled = true
			break STORIES
		default:
		}

	}

	wg.Wait()

	if nSuccess != int64(n) {
		if contextCanceled {
			return results, ctx.Err()
		}
		return results, fmt.Errorf("Didn't successfully fetch all items.")
	}

	return results, nil
}
