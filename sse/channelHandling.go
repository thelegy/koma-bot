package sse

import (
	"container/list"
	"time"
)

type internalClient struct {
	channel chan<- interface{}
	quit    <-chan struct{}
	strikes int
}

// Create new clients and send them out via the channel
func createChannels(c chan<- Client, clients *list.List, done <-chan struct{}) {
	for {
		newClientChannel := make(chan interface{})
		newQuitChannel := make(chan struct{})
		newClient := Client{newClientChannel, newQuitChannel}
		newInternalClient := internalClient{newClientChannel, newQuitChannel, 0}

		select {
		case c <- newClient:
			clients.PushBack(&newInternalClient)
		case <-done:
			close(c)
			// let gc do something for its money
			return
		}
	}
}

// Serve the clients their events and destroy them when its time
// This also triggers the done chan, when c is closed
func serveChannels(c <-chan interface{}, clients *list.List, done chan<- struct{}) {
	for {
		var (
			event       interface{}
			channelOpen bool
		)
		event = nil
		channelOpen = true
		select {
		case event, channelOpen = <-c:
		case <-time.After(120 * time.Second):
		}

		if !channelOpen {
			close(done)

			// wait so no new clients come any longer
			<-time.After(1 * time.Second)
			for client := range clientsIter(clients) {
				close(client.value.channel)
			}

			// empty all references
			clients.Init()
			return
		}

		// channel is still open
		for client := range clientsIter(clients) {
			select {
			case <-client.value.quit:
				// client has quit
				client.removeFrom(clients)
				continue
			default:
			}

			// client has not quit
			if event != nil {
				select {
				case client.value.channel <- event:
				case <-time.After(100 * time.Millisecond):
					// strike
					client.value.strikes++
					if client.value.strikes >= 5 {
						client.removeFrom(clients)
					}
				}
			}
		}
	}
}
