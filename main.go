package main

import (
	"log"
	"net"
	"os"
)

/* global vars */
var (
	events    = initEvents()
	clients   = initClients()
	followers = make(map[int]map[int]int)
	seqNum    = 1
	noEvents  = false
)

func main() {

	/* start the server */
	server, e := Start()
	log.Printf("Server started...\n")

	if e != nil {
		close(server.Done)
		os.Exit(1)
	}

	/* start processing */
	go run(server.eventListener, server.Done, AcceptEventSource)
	go run(server.clientListener, server.Done, AcceptClients)

	/* start sending */
	sendMessages(server.Done)

}

type acceptFunc func(l net.Listener)

/* runner for listeners */
func run(l net.Listener, done chan struct{}, f acceptFunc) {
	for {
		go f(l)

		select {
		case <-done:
			_ = l.Close()
			return
		}
	}
}

/* sends events in ascending order */
func sendMessages(done chan struct{}) {
	for {
		if noEvents && events.isEmpty() {
			flushAll()
			close(done)
			return
		} else if e, ok := events.Get(seqNum); ok {
			events.Delete(seqNum)
			handleEvent(e)
			seqNum++
		}
	}
}
