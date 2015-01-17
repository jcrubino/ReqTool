package main

import (
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
	Domain   []byte // tells Urlscanner the domain for dedicated connection
	GetUrls  []byte // redis key of all URLs to Get
	PostUrls []byte // redis key of all URLs to Post to
	Sleep    int    // milliseconds between goroutine execution
	TimeOut  int    // milliseconds http client timeout
	Pause    int    // milliseconds pause between request sets on one domain
	Run      bool   // state of execution
	Msg      []byte // use msgs for rpc
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

const FULLSTOP = Command{[]byte("*"), []byte(""), []byte(""), 10000, 10000, 10000, false, []byte("FULLSTOP")}
