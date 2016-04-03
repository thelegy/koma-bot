package sse

import (
	"container/list"
	"time"
)

type iterClient struct {
	element *list.Element
	value   *internalClient
}

func (client iterClient) removeFrom(clients *list.List) {
	close(client.value.channel)
	clients.Remove(client.element)
}

func clientsIter(clients *list.List) <-chan iterClient {
	len := clients.Len()
	c := make(chan iterClient, len)
	go func(clients *list.List, c chan<- iterClient) {
		for e := clients.Front(); e != nil; e = e.Next() {
			client, ok := e.Value.(*internalClient)
			if !ok {
				continue
			}
			clientContainer := iterClient{e, client}
			select {
			case c <- clientContainer:
			case <-time.After(5 * time.Second):
				// iteration seems to be aborted
				close(c)
				return
			}
		}
		close(c)
	}(clients, c)
	return c
}
