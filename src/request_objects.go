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
	Domain   []byte // tells Urlscanner the domain to work on
	GetUrls  []byte // redis key of all URLs to Get
	PostUrls []byte // redis key of all URLs to Post to
	Sleep    int    // length between goroutine execution
	TimeOut  int    // http client timeout
	Run      bool   // state of execution
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

const FULLSTOP = Command{[]byte("*"), []byte(""), []byte(""), 1000, 1000, false}
