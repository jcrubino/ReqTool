package main

import (
	"fmt"
)

// function to AutoADjust Network Requests
// http.client TimeOut and throttle between goroutines
func TuneRequests() {
	var counter int = 0
	var start int64 = 0
	var mps float64
	for msg := range results {

		if counter%5 == 0 {
			mps = float64(counter) / float64(msg.Time.Unix()-start)
			fmt.Println(msg.Key, msg.Err, mps)
		}
		if counter == 0 {
			start = int64(msg.Time.Unix())
		}
		if msg.Err == nil {
			counter++
		}

	}
}
