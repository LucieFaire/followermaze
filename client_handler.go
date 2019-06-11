package main

import (
	"bufio"
	"log"
	"net"
	"strconv"
	"sync"
)

type ClientHandler struct {
	id     int
	conn   net.Conn
	io     *sync.Mutex
	writer *bufio.Writer
	fMap   FollowMap
}

func initClientHandler(conn net.Conn) *ClientHandler {
	return &ClientHandler{nil, conn, &sync.Mutex{}, nil, initFollowers()}
}

func Setup(conn net.Conn) {
	handler := initClientHandler(conn)

	defer conn.Close()
	defer clients.Delete(handler.id)

	handler.read()
}

func (h *ClientHandler) Write(e *event) {
	_, err := h.writer.WriteString(e.msg)
	if err != nil {
		log.Panicf("Could not send a message to a client: %s", err.Error())
	}
}

func (h *ClientHandler) read() {
	reader := bufio.NewReader(h.conn)
	for {
		in, err := reader.ReadString(lineDelimiter)
		if err != nil {
			log.Fatalf("Could not read client's input: %s", err)
			return
		}
		id, err := strconv.Atoi(in)
		if err != nil {
			log.Fatalf("Could not extract the client id: %s", err)
			return
		}
		h.id = id
		h.writer = bufio.NewWriter(h.conn)
		clients.Store(h.id, h)
	}
}

func (h *ClientHandler) Follow(e *event) {
	if v, ok := clients.Load(e.from); ok {
		h := v.(*ClientHandler)
		h.fMap.Put(h)
		h.Write(e)
	}
}

func (h *ClientHandler) UnFollow(e *event) {
	h.fMap.Delete(e.from)

}

func (h *ClientHandler) UpdateSt(e *event) {
	for _, f := range h.fMap.followers {
		f.Write(e)
	}
}

/* synchronized map with fMap */
type FollowMap struct {
	followers map[int]*ClientHandler
	m         sync.RWMutex
}

func initFollowers() FollowMap {
	return FollowMap{
		make(map[int]*ClientHandler),
		sync.RWMutex{},
	}
}

func (f *FollowMap) Put(h *ClientHandler) {
	f.m.Lock()
	defer f.m.Unlock()

	f.followers[h.id] = h
}

func (f *FollowMap) Get(id int) *ClientHandler {
	f.m.RLock()
	defer f.m.RUnlock()

	if h, ok := f.followers[id]; ok {
		return h
	}
	return nil
}

func (f *FollowMap) Delete(id int) {
	f.m.Lock()
	defer f.m.Unlock()

	delete(f.followers, id)
}

func (f *FollowMap) Contains(id int) bool {
	f.m.RLock()
	defer f.m.RUnlock()

	_, ok := f.followers[id]
	return ok
}
