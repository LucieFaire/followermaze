package main

import (
	"net"
	"os"
)

func main() {

	//start the server
	server, e := Start()

	if e != nil {
		GracefullyClose(server)
		os.Exit(1)
	}

	go runE(server.eventListener, server.Done)
	go runCli(server.clientListener, server.Done)

	//accept the event source connection
	AcceptEventSource(server.clientListener)

}

func runE(l net.Listener, done chan struct{}) {

}

func runCli(l net.Listener, done chan struct{}) {
	for {
		go AcceptClients(l)

		select {
		case <-done:
			_ = l.Close()
			return
			
		}
	}
}
