package main

import (
	"github.com/labstack/gommon/log"
	"net"
	"os"
	"sync"
)

var (
	// init the event map
	events  = make(map[int]*event)
	clients = sync.Map{}
	seqNum  = 1
)

func main() {

	/* start the server */
	server, e := Start()

	if e != nil {
		GracefullyClose(server)
		os.Exit(1)
	}

	/* start processing */
	go run(server.eventListener, server.Done, AcceptEventSource)
	go run(server.clientListener, server.Done, AcceptClients)

	/* start sending */
	go sendMessages(server.Done)

}

type acceptFunc func(l net.Listener)

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

func sendMessages(done chan struct{}) {
	for {
		if seqNum == 10000000 {
			close(done)
		} else if e, ok := events[seqNum]; ok {
			delete(events, seqNum)
			handleEvent(e)
			seqNum++
		} else {
			break
		}
	}
}

func handleEvent(e *event) {
	switch e.eType {
	case Follow:
		if v, ok := clients.Load(e.to); ok {
			h := v.(*ClientHandler)
			h.Follow(e)
		}
	case Unfollow:
		if v, ok := clients.Load(e.to); ok {
			h := v.(*ClientHandler)
			h.UnFollow(e)
		}
	case Broadcast:
		clients.Range(func(_, v interface{}) bool {
			h := v.(*ClientHandler)
			go h.Write(e)
			return true
		})
	case PrivateMessage:
		if v, ok := clients.Load(e.to); ok {
			h := v.(*ClientHandler)
			h.Write(e)
		}
	case StatusUpdate:
		if v, ok := clients.Load(e.to); ok {
			h := v.(*ClientHandler)
			h.UpdateSt(e)
		}
	default:
		log.Fatalf("Could not recognize event type %v", e.eType)
	}
}
