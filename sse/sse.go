package sse

import (
	"container/list"
)

// Control Frontend to the SSE logic
type Provider struct {
	NewClients  <-chan Client      // chan yielding new clients
	EventStream chan<- interface{} // chan to send events to
}

// SSE Client retreiving events
type Client struct {
	Channel <-chan interface{} // chan yielding the events
	Quit    chan<- struct{}    // chan to signal when client is gone
}

// Creates a new Provider and starts up the logic
func NewProvider() Provider {
	done := make(chan struct{})
	clients := list.New()

	newClients := make(chan Client)
	eventStream := make(chan interface{})

	eventProvider := Provider{newClients, eventStream}

	// Create new clients and send them out via newClients
	go createChannels(newClients, clients, done)

	// Serve channels including closing them when its time
	go serveChannels(eventStream, clients, done)

	return eventProvider
}
