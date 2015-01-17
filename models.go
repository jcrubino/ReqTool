package main

import (
	"fmt"
	"time"
)

//  Result msg per request
type Result struct {
	Key   string
	Value []byte
	Err   error
	Time  time.Time
}

// tuning commands to GoRoutine Gatling
type Command struct {
	Sleep   int
	TimeOut int
	Run     bool
	Idle    int // minutes
}

// Config for daemon
type Config struct {
	Headers        string
	PIDName        string
	GetUrls        string
	PostUrls       string
	WSUrl          string
	WSPort         string
	AuthID         string
	AuthSecret     string
	AuthScheme     string
	UpdateUnits    int
	UpdateInterval int
}

// shutdown all requests
func shutdown() {
	control <- Command{1, 2, false, 0}
	fmt.Println("Shutdown command sent")
	close(results)
	close(control)
}

// idle requests
func idle(ms int) {
	control <- Command{1, 2, true, ms}
	fmt.Printf("Idling command sent: %v minutes\n", ms)

}

var results = make(chan Result, 120)
var control = make(chan Command)
