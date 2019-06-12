package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCorrect_ProcessInput_Follow(t *testing.T) {
	processInput("666|F|60|50\n")

	e, ok := events.Get(666)

	assert.True(t, ok)
	assert.EqualValues(t, &event{666, Follow, 60, 50, "666|F|60|50\n"}, e)
}

func TestCorrect_ProcessInput_Broadcast(t *testing.T) {
	processInput("54|B\n")

	e, ok := events.Get(54)

	assert.True(t, ok)
	assert.EqualValues(t, &event{54, Broadcast, 0, 0, "54|B\n"}, e)
}

func TestCorrect_ProcessInput_Status(t *testing.T) {
	processInput("634|S|32\n")

	e, ok := events.Get(634)

	assert.True(t, ok)
	assert.EqualValues(t, &event{634, StatusUpdate, 32, 0, "634|S|32\n"}, e)
}

func TestFail_ProcessInput(t *testing.T) {
	processInput("634|\n")

	t.Log("Unknown event type, could not parse the input\n")
}
