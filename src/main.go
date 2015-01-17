package main

import (
	"fmt"
	"github.com/takama/daemon"
	"os"
)

var results = make(chan Result, 120)
var control = make(chan Command)

func main() {

	// dameon init
	srv, err := daemon.New(name, description)
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
	service := &Service{srv}
	status, err := service.Manage()
	if err != nil {
		fmt.Println(status, "\nError: ", err)
		os.Exit(1)
	}
	fmt.Println(status)

	defer shutdown()

}
