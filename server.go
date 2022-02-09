package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

var count int = 0
var mtx_lock sync.RWMutex

func rootCaller(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Invalid path!!")
	fmt.Fprintf(w, "Usage: http://localhost:8080/items/<item_id>")
}

// Server code to retrieve items; for debugging purpose only
func itemCaller(w http.ResponseWriter, r *http.Request) {
	_id := r.URL.Path[7:]
	auth := r.Header.Get("Authorization")
	decodeAuth, err := base64.StdEncoding.DecodeString(auth)
	if err != nil || string(decodeAuth) != _id {
		fmt.Fprintf(w, "Invalid authorization: Authorization header mtx_lockst be base 64 encoded ID")
	} else {
		fmt.Println(r.Method, r.URL.Path, r.RemoteAddr, r.UserAgent())

		// mtx_locktex so that count is not accessed/updated by conc. req. handlers at the same time
		mtx_lock.Lock()
		count++

		// If number of concurrent requests exceeds 5, return error 429
		if count > 5 {
			mtx_lock.Unlock()
			w.WriteHeader(429)
			fmt.Fprintf(w, "ERROR 429: Too many concurrent requests made")
		} else if rand.Intn(41) == 4 {
			// Simulate error generation randomly
			mtx_lock.Unlock()
			w.WriteHeader(404)
			fmt.Fprintf(w, "ERROR 404")
		} else {
			mtx_lock.Unlock()
			fmt.Fprintf(w, "%d", rand.Intn(100000000000000000))
		}

		// Simulation of delay in processing the request
		time.Sleep(time.Second)

		mtx_lock.Lock()
		count--
		mtx_lock.Unlock()
	}
}

func main() {
	http.HandleFunc("/", rootCaller)
	http.HandleFunc("/items/", itemCaller)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
