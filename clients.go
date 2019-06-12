package main

import (
	"sync"
)

/* synchronized map with connected clients */
type Clients struct {
	set map[int]*ClientHandler
	m   sync.RWMutex
}

func initClients() *Clients {
	return &Clients{
		make(map[int]*ClientHandler),
		sync.RWMutex{},
	}
}

func (f *Clients) Put(h *ClientHandler) {
	f.m.Lock()
	defer f.m.Unlock()

	f.set[h.id] = h
}

func (f *Clients) Get(id int) (*ClientHandler, bool) {
	f.m.RLock()
	defer f.m.RUnlock()
	h, ok := f.set[id]
	return h, ok
}

func (f *Clients) Delete(id int) {
	f.m.Lock()
	defer f.m.Unlock()

	delete(f.set, id)
}

func (f *Clients) Contains(id int) bool {
	f.m.RLock()
	defer f.m.RUnlock()

	_, ok := f.set[id]
	return ok
}

type rgFunc func(h *ClientHandler)

func (f *Clients) Range(someFunc rgFunc) {
	f.m.Lock()
	defer f.m.Unlock()

	for _, v := range f.set {
		someFunc(v)
	}
}

func (f *Clients) isEmpty() bool {
	f.m.RLock()
	defer f.m.RUnlock()

	return len(f.set) == 0
}
