package main

import "sync"

/* synchronized map with ordered events */
type Events struct {
	set map[int]*event
	m   sync.RWMutex
}

func initEvents() *Events {
	return &Events{
		make(map[int]*event),
		sync.RWMutex{},
	}
}

func (s *Events) Put(e *event) {
	s.m.Lock()
	defer s.m.Unlock()

	s.set[e.seqId] = e
}

func (s *Events) Get(id int) (*event, bool) {
	s.m.RLock()
	defer s.m.RUnlock()

	h, ok := s.set[id]
	return h, ok

}

func (s *Events) Delete(id int) {
	s.m.Lock()
	defer s.m.Unlock()

	delete(s.set, id)
}

func (s *Events) Contains(id int) bool {
	s.m.RLock()
	defer s.m.RUnlock()

	_, ok := s.set[id]
	return ok
}

func (s *Events) isEmpty() bool {
	s.m.RLock()
	defer s.m.RUnlock()

	return len(s.set) == 0
}

type eFunc func(e *event)

func (s *Events) Range(someFunc eFunc) {
	s.m.Lock()
	defer s.m.Unlock()

	for _, v := range s.set {
		someFunc(v)
	}
}
