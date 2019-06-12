package main

import (
	"log"
	"net"
)

const (
	clientAddr = ":9099"
	eventAddr  = ":9090"
)

/* Server structure with multiple listeners for clients and events and priority queue for incoming events */
type Server struct {
	clientListener net.Listener
	eventListener  net.Listener
	Done           chan struct{}
}

/* Starts the server */
func Start() (*Server, error) {

	// start the server to listen to events
	eventListener, e := net.Listen("tcp", eventAddr)
	if e != nil {
		log.Fatalf("Could not start the server on port: %s. %s", eventAddr, e.Error())
		return nil, e
	}

	// start the server to listen to client connections
	clientListener, e := net.Listen("tcp", clientAddr)
	if e != nil {
		log.Fatalf("Could not start the server on port: %s. %s", clientAddr, e.Error())
		return nil, e
	}

	Done := make(chan struct{})

	return &Server{clientListener, eventListener, Done}, nil
}

/* Processes the incoming connection from client */
func AcceptClients(l net.Listener) {
	for {
		conn, e := l.Accept()
		if e != nil {
			log.Printf("Failed to accept a client connection with error: %s\n", e.Error())
			return
		}
		go Setup(conn)
	}
}

/* Processes the event source connection */
func AcceptEventSource(l net.Listener) {
	conn, e := l.Accept()

	defer signal()
	defer l.Close()
	defer conn.Close()

	if e != nil {
		log.Printf("Failed to accept an event source connection with error: %s\n", e.Error())
		return
	}
	handler := InitEventHandler(conn)
	handler.read()
}

/* signals to stop the server */
func signal() {
	noEvents = true
}

/* handles events according to the type */
func handleEvent(e *event) {
	switch e.eType {
	case Follow:
		follow(e)
	case Unfollow:
		unFollow(e)
	case Broadcast:
		clients.Range(func(h *ClientHandler) {
			h.Write(e)
		})
	case PrivateMessage:
		if h, ok := clients.Get(e.to); ok {
			h.Write(e)
		}
	case StatusUpdate:
		updateStatus(e)
	default:
		log.Fatalf("Could not recognize event type %v", e.eType)
	}
}

/* Logic for status update event */
func updateStatus(e *event) {
	fMap := followers[e.from]
	for _, f := range fMap {
		if h, ok := clients.Get(f); ok {
			h.Write(e)
		}
	}
}

/* Logic for follow event */
func follow(e *event) {
	fMap, ok := followers[e.to]
	if !ok {
		fMap = make(map[int]int)
	}
	fMap[e.from] = e.from
	followers[e.to] = fMap
	if h, ok := clients.Get(e.to); ok {
		h.Write(e)
	}
}

/* Logic for unfollow event */
func unFollow(e *event) {
	if fMap, ok := followers[e.to]; ok {
		delete(fMap, e.from)
	}
}
