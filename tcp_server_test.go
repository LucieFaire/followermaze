package main

import (
	"github.com/stretchr/testify/assert"
	"net"
	"sync"
	"testing"
)

func TestCorrect_InitialFollow(t *testing.T) {
	fakeHandler()
	e := &event{666, Follow, 60, 50, "666|F|60|50\n"}

	follow(e)

	fMap, ok := followers[50]

	assert.True(t, ok)
	assert.Equal(t, 1, len(fMap))
	assert.Equal(t, 60, fMap[60])
}

func TestCorrect_NextFollow(t *testing.T) {
	fakeHandler()
	fakeFollowers(22, 50)
	e := &event{666, Follow, 60, 50, "666|F|60|50\n"}

	follow(e)

	fMap, ok := followers[50]

	assert.True(t, ok)
	assert.Equal(t, 2, len(fMap))
	assert.Equal(t, 60, fMap[60])
}

func TestCorrect_Unfollow(t *testing.T) {
	fakeFollowers(60, 50)
	e := &event{666, Unfollow, 60, 50, "666|U|60|50\n"}

	unFollow(e)

	fMap, ok := followers[50]

	assert.True(t, ok)
	assert.Equal(t, 0, len(fMap))
}

func fakeHandler() {
	conn, _ := net.Pipe()
	handler := &ClientHandler{-1, conn, &sync.Mutex{}}

	clients.Put(handler)
}

func fakeFollowers(to int, from int) {
	fMap := make(map[int]int)
	fMap[to] = to
	followers[from] = fMap
}
