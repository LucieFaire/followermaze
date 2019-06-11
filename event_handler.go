package main

import (
	"bufio"
	"github.com/labstack/gommon/log"
	"net"
	"strconv"
	"strings"
	"sync"
)

const (
	lineDelimiter    = '\n'
	payloadDelimiter = "|"
	Follow           = iota
	Unfollow
	Broadcast
	PrivateMessage
	StatusUpdate
)

type event struct {
	seqId int
	eType int
	from  int
	to    int
	msg string
}

type EventHandler struct {
	conn net.Conn
	events sync.Map
}

func InitEventHandler(conn net.Conn) *EventHandler {
	return &EventHandler{conn, sync.Map{}}
}

func (h *EventHandler) read() {
	reader := bufio.NewReader(h.conn)
	for {
		in, err := reader.ReadString(lineDelimiter)
		if err != nil {
			log.Fatalf("Could not read event: %s", err)
			return
		}
		h.processInput(in)
		go h.sendEvents()
	}
}

func (h *EventHandler) sendEvents() {

}

func (h *EventHandler) processInput(s string) {
	params := strings.Split(s, payloadDelimiter)
	if len(params) < 2 {
		log.Fatalf("Error while parsing the input: wrong format")
		return
	}

	e := &event{}
	e.seqId, _ = strconv.Atoi(params[0])
	e.msg = s

	switch params[1] {
	case "F":
		e.eType = Follow
		setToFrom(e, params)
	case "U":
		e.eType = Unfollow
		setToFrom(e, params)
	case "B":
		e.eType = Broadcast
	case "P":
		e.eType = PrivateMessage
		setToFrom(e, params)
	case "S":
		e.eType = StatusUpdate
		e.from, _ = strconv.Atoi(params[2])
	default:

	}
	h.events.Store(e.seqId, e)
}

func setToFrom(e *event, params []string) {
	e.from, _ = strconv.Atoi(params[2])
	e.to, _ = strconv.Atoi(params[3])
}
