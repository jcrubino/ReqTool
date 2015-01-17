package main

import (
	"bytes"
	"fmt"
	"github.com/fzzy/radix/extra/pool"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

// Gets a connection from the redis pool
func SendToRedis(p *pool.Pool, key string, value []byte) {
	conn, redisErr := p.Get()
	if redisErr != nil {
		fmt.Println("Redis Pool Error: ", redisErr)

	}
	defer p.CarefullyPut(conn, &redisErr)
	// publish to redis
	redisErr = conn.Cmd("publish", key, value).Err
	if redisErr != nil {
		fmt.Println("Redis Publish Error: ", redisErr)

	}

	if !bytes.Equal(value, []byte("NA")) {
		// set to redis
		redisErr = conn.Cmd("SET", key, value).Err
		if redisErr != nil {
			fmt.Println("Redis Set Key Error: ", redisErr)

		}
		// expire redis key after 10 seconds
		redisErr = conn.Cmd("EXPIRE", key, 10).Err
		if redisErr != nil {
			fmt.Println("Redis Expire Error: ", redisErr)
		}
	}
}

func URLScanner(results chan Result, command chan Command, setup Command) {
	fmt.Println("Urlscanner Starting")
	// tuneable params

	domain := setup.Domain
	getUrls := setup.GetUrls   // key in redis to get GET urls
	postUrls := setup.PostUrls // key in redis to get POST urls
	throttle := setup.Sleep
	pause := setup.Pause
	timeout := time.Duration(setup.TimeOut) * time.Millisecond

	p, err := pool.NewPool("tcp", "127.0.0.1:6379", 25)
	client := &http.Client{Timeout: timeout}
	if err != nil {
		fmt.Println("Redis Connection Error")
		return
	}

	for {

		// Pause per set of requests
		time.Sleep(time.Millisecond * time.Duration(pause))
		select {

		case msg := <-command:
			if bytes.Equal(msg.Domain, domain) || bytes.Equal(msg.Domain, []byte("*")) {

				if msg.Run == true {
					TIMEOUT := msg.TimeOut
					SLEEP := msg.Sleep
					timeout := time.Duration(TIMEOUT) * time.Millisecond

				}
				if msg.Run == false {
					return
				}
			}
		default:

			var wg sync.WaitGroup
			wg.Add(len(getUrls))

			for k, url := range getUrls {

				// throttle per goroutine
				time.Sleep(time.Millisecond * time.Duration(throttle))

				// select / case / match
				// check for messages that might be to post data
				// post execution get priority over next goroutine

				go func(k string, url string, wg *sync.WaitGroup, headers map[string]string) {
					req, err := http.NewRequest("GET", url, nil)
					for k, v := range headers {
						req.Header.Add(k, v)
					}
					resp, err := client.Do(req)
					timeStamp := time.Now()
					if err != nil {
						SendToRedis(p, k, []byte(""))
						results <- Result{k, []byte(""), err, timeStamp}
						wg.Done()
						return
					}
					defer resp.Body.Close()
					bs, err := ioutil.ReadAll(resp.Body)
					if err != nil {
						SendToRedis(p, k, []byte(""))
						results <- Result{k, []byte(""), err, timeStamp}
						wg.Done()
						return
					}
					SendToRedis(p, k, bs)
					results <- Result{k, bs, nil, timeStamp}
					wg.Done()
					return

				}(k, url, &wg)
			}

		}
	}
}

// shutdown all requests
func shutdown() {
	control <- FULLSTOP
	fmt.Println("Shutdown command sent")
	close(results)
	close(control)
}
