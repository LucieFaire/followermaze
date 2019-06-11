package main

import (
	"net"
	"sync"
)

type ClientHandler struct {
	id        int
	conn      net.Conn
	io        *sync.Mutex
	followers FollowMap
}

func initClientHandler(conn net.Conn) *ClientHandler {
	return &ClientHandler{nil, conn, &sync.Mutex{}, initFollowers()}
}

func Setup(conn net.Conn) {
	handler := initClientHandler(conn)

}

/* synchronized map with followers */
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
