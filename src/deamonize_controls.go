package main

import (
	//"github.com/BurntSushi/toml"
	"github.com/takama/daemon"
)

// service object acted upon by daemon lib
type Service struct {
	daemon.Daemon
}

func (service *Service) Start() (string, error) {

	// init toml file or redis

	go URLScanner(results, control, urls)
	return "", nil
}

func (service *Service) Stop() (string, error) {
	// stop domain
	// stop ReqTool

	control <- FULLSTOP
	return "", nil
}

func (service *Service) Status() (string, error) {
	// pass
	return "", nil
}

func (service *Service) Update() (string, error) {
	// pass
	return "", nil
}
