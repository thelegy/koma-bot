package sse

import (
	"testing"

	"container/list"
)

func TestClientsIter(t *testing.T) {
	clients := list.New()
	const test_num = 15

	for i := 0; i < test_num; i++ {
		newClientChannel := make(chan interface{})
		newQuitChannel := make(chan struct{})
		newInternalClient := internalClient{newClientChannel, newQuitChannel, 0}

		clients.PushBack(&newInternalClient)
	}

	k := 0

	for _ = range clientsIter(clients) {
		k++
	}

	if k < test_num {
		t.Error("Iterated over not enough items")
	}

	if k > test_num {
		t.Error("Iterated over too many items")
	}
}
