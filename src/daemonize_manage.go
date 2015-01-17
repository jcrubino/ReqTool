package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

const (
	name        = "ReqTool"
	description = "A Concurrent HTTP Client and Publisher"
	port        = ":9977"
)

func (service *Service) Manage() (string, error) {
	usage := "Usage: reqtool install | remove | start | stop | status "
	if len(os.Args) > 1 {
		command := os.Args[1]
		switch command {
		case "install":
			return service.Install()
		case "remove":
			return service.Remove()
		case "start":
			return service.Start()
		case "stop":
			return service.Stop()
		case "status":
			return service.Status()
		default:
			return usage, nil
		}
	}

	// Listener Channel
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, os.Kill, syscall.SIGTERM)

	listener, err := net.Listen("tcp", port)
	if err != nil {
		return "Possible port binding problem.", err
	}

	// Connection Accept Channel
	listen := make(chan net.Conn, 100)
	go acceptConnection(listener, listen)

	// Listen Handle Loop
	for {
		select {
		case conn := <-listen:
			go handleClient(conn)
		case killSignal := <-interrupt:
			log.Println("Got signal:", killSignal)
			log.Println("Stoping listening on ", listener.Addr())
			listener.Close()
			if killSignal == os.Interrupt {
				return "Daemon was interruped by system signal", nil
			}
			return "Daemon was killed", nil
		}
	}

	return usage, nil
}

func acceptConnection(listener net.Listener, listen chan<- net.Conn) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		listen <- conn
	}
}

func handleClient(client net.Conn) {
	for {
		buf := make([]byte, 4096)
		numbytes, err := client.Read(buf)
		if numbytes == 0 || err != nil {
			return
		}
		client.Write(buf)
	}
}
