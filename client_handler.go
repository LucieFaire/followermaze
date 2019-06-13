package main

import (
	"bufio"
	"log"
	"net"
	"strconv"
	"strings"
	"sync"
)

/* handler per each client connection */
type ClientHandler struct {
	id     int
	conn   net.Conn
	mutex  *sync.Mutex
	writer *bufio.Writer
}

func initClientHandler(conn net.Conn) *ClientHandler {
	return &ClientHandler{-1, conn, &sync.Mutex{}, nil}
}

/* go routine with client setup and management */
func Setup(conn net.Conn) {
	handler := initClientHandler(conn)

	defer conn.Close()
	defer clients.Delete(handler.id)

	handler.read()
}

/* writes to the client connection */
func (h *ClientHandler) Write(e *event) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	_, err := h.writer.WriteString(e.msg)
	if err != nil {
		log.Printf("Could not send a message to a client: %s\n", err.Error())
	}
}

const buffSize = 2048

/* reads the client id and checks the connection is open */
func (h *ClientHandler) read() {
	reader := bufio.NewReader(h.conn)
	for {
		in, err := reader.ReadString(lineDelimiter)
		if err != nil {
			log.Printf("No input: %s\n", err)
			return
		}
		id, err := strconv.Atoi(strings.TrimRight(in, string(lineDelimiter)))
		if err != nil {
			log.Printf("Could not extract the client id: %s\n", err)
			return
		}
		h.id = id
		h.writer = bufio.NewWriterSize(h.conn, buffSize)
		clients.Put(h)
	}
}

func (h *ClientHandler) Flush() {
	_ = h.writer.Flush()
}
