package main

import (
	"bufio"
	"github.com/labstack/gommon/log"
	"net"
	"strconv"
	"strings"
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
	eType uint
	from  int
	to    int
	msg   string
}

type EventHandler struct {
	conn net.Conn
}

func InitEventHandler(conn net.Conn) *EventHandler {
	return &EventHandler{conn}
}

func (h *EventHandler) read() {
	reader := bufio.NewReader(h.conn)
	for {
		in, err := reader.ReadString(lineDelimiter)
		if err != nil {
			log.Fatalf("Could not read event: %s", err)
			return
		}
		go processInput(in)
	}
}

func processInput(s string) {
	params := strings.Split(s, payloadDelimiter)
	if len(params) < 2 {
		log.Fatalf("Error while parsing the input: wrong format\n")
		return
	}

	e := &event{}
	e.seqId, _ = strconv.Atoi(params[0])
	e.msg = s

	switch params[1] {
	case "F":
		e.eType = Follow
		setFromTo(e, params)
	case "U":
		e.eType = Unfollow
		setFromTo(e, params)
	case "B":
		e.eType = Broadcast
	case "P":
		e.eType = PrivateMessage
		setFromTo(e, params)
	case "S":
		e.eType = StatusUpdate
		e.from, _ = strconv.Atoi(params[2])
	default:
		log.Fatal("Unknown event type, could not parse the input\n")
	}
	events.Store(e.seqId, e)
}

func setFromTo(e *event, params []string) {
	e.from, _ = strconv.Atoi(params[2])
	e.to, _ = strconv.Atoi(params[3])
}
