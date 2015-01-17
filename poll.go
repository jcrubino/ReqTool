package main

/*
   Daemonize https://github.com/takama/daemon
*/
import (
	"time"
)

func main() {

	go UrlScanner(results, control, Urls)
	time.Sleep(time.Duration(90) * time.Second)
	defer shutdown()

}
