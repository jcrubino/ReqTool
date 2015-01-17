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

func UrlScanner(results chan Result, command chan Command, urls map[string]string) {
	fmt.Println("Urlscanner Starting")
	// tuneable params
	var TIMEOUT int = 575
	var SLEEP int = 30
	timeout := time.Duration(TIMEOUT) * time.Millisecond
	p, err := pool.NewPool("tcp", "127.0.0.1:6379", 25)
	client := &http.Client{Timeout: timeout}
	if err != nil {
		fmt.Println("Redis Connection Error")
		return
	}

	for {
		select {

		case cmd := <-command:

			if cmd.Run == true {
				TIMEOUT = cmd.TimeOut
				SLEEP = cmd.Sleep
				timeout = time.Duration(TIMEOUT) * time.Millisecond
				fmt.Println("Setting new Command")
				time.Sleep(time.Minute * time.Duration(cmd.Idle))

			} else {
				fmt.Println("Shutting Down")
				break
			}
		default:

			// Throttle:  Change to AutoTune Function
			time.Sleep(time.Millisecond * time.Duration(1))

			var wg sync.WaitGroup
			wg.Add(len(Urls))

			for k, url := range urls {
				// ToDo Throttle: Change to AutoTune Function
				time.Sleep(time.Millisecond * time.Duration(SLEEP))

				// select / case / match
				// check for messages that might be to post data
				// trade execution get priority over next goroutine
				// reuses the same connection to server

				go func(k string, url string, wg *sync.WaitGroup) {
					req, err := http.NewRequest("GET", url, nil)
					req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
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
