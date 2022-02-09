package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

// Time intervals for backoff
var max_bkoff int = 4
var bkoff_times = []time.Duration{
	1 * time.Second,
	2 * time.Second,
	4 * time.Second,
	8 * time.Second,
}

// Waitgroup to hold till all items are queried
var wait_gp = &sync.WaitGroup{}

// Dictionary that holds items already queried and corresponding mutex
var dict map[string]bool = make(map[string]bool)
var mtx_lock sync.RWMutex

// Backoff routine to add the failed ID back to the channel after timeout
func backoffHandler(item string, iter int, channel chan []interface{}) {
	// Wait for timer to expire for retries, then push back to channel and increment timer value
	time.Sleep(bkoff_times[iter])
	channel <- []interface{}{item, iter + 1}
}

// Return response for GET request
func getResponse(base_url string, _id string) []interface{} {
	// Convert id to base64
	auth := base64.URLEncoding.EncodeToString([]byte(_id))
	path := base_url + _id

	cli := &http.Client{}
	req, _ := http.NewRequest("GET", path, nil)
	req.Header.Set("Authorization", auth)

	// GET request sent and error, if any, displayed
	res, err := cli.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	// Read GET request's response
	resp, _ := ioutil.ReadAll(res.Body)
	code := res.StatusCode
	res.Body.Close()

	// Return response
	return []interface{}{string(resp[:]), _id, code}
}

// Main routine to getResponse() of GET request inorder for the IDs received from the channel
func requestInfo(base_url string, channel chan []interface{}) {
	// Iterate over all query items in the channel
	for item := range channel {
		resp := getResponse(base_url, item[0].(string))

		if resp[2].(int) != 200 {
			// Backoff routine to be called if return code is not 200
			if item[1].(int) < max_bkoff {
				go backoffHandler(item[0].(string), item[1].(int), channel)
			} else {
				wait_gp.Done()

				mtx_lock.Lock()
				dict[item[0].(string)] = false
				mtx_lock.Unlock()
			}
		} else {
			// Wait group processing is completed successfully for this item
			wait_gp.Done()
			fmt.Printf("Item[ %s ] returned information [ %s ]\n", resp[1].(string), resp[0].(string))
		}
	}
}

func main() {
	// base_url := "http://localhost:8080/items/"   // Testing with server.go
	base_url := "https://challenges.qluv.io/items/" // Actual server to query

	// Randomly generate a list of items for testing
	item_list := make([]string, 1000)
	for i := 0; i < len(item_list); i++ {
		item_list[i] = fmt.Sprint(rand.Intn(380))
	}
	// make a channel with a capacity to hold as many items
	channel := make(chan []interface{}, 5)

	// to measure time taken for all the queries
	start := time.Now()

	/* Create 5 Go routines to keep sending GET requests simultaneously and keep sending requests to
	these routines using a channel */
	for idx := 0; idx < 5; idx++ {
		go requestInfo(base_url, channel)
	}

	// Process items not queried yet by pushing their IDs to channel for processing
	for current := 0; current < len(item_list); current++ {
		if !dict[item_list[current]] {
			mtx_lock.Lock()
			dict[item_list[current]] = true
			mtx_lock.Unlock()

			wait_gp.Add(1)
			channel <- []interface{}{item_list[current], 0}
		}
	}

	// Wait for all the items in channel to be processed
	wait_gp.Wait()

	close(channel)

	duration := time.Since(start)
	fmt.Println("Total time taken = ", duration)
}
