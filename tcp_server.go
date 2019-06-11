package main

import (
	"github.com/labstack/gommon/log"
	"net"
)

const (
	clientAddr = "9099"
	eventAddr  = "9090"
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

type Handler func()

/* Processes the incoming connection from client */
func AcceptClients(l net.Listener) {
	for {
		conn, e := l.Accept()
		if e != nil {
			log.Printf("Failed to accept a client connection with error: %s", e.Error())
			return
		}
		go Setup(conn)

	}
}

func GracefullyClose(s *Server) {
	log.Fatal("Exiting...")
	close(s.Done)
}
