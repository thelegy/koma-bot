package sse

import (
	"testing"
	"time"
)

func TestNewProvider(t *testing.T) {
	provider := NewProvider()

	_ = provider
}

func TestShutdown(t *testing.T) {
	provider := NewProvider()

	close(provider.EventStream)
}

func TestGetClients(t *testing.T) {
	provider := NewProvider()
	const clients_num = 1000
	var clients [clients_num]Client

	for i := 0; i < clients_num; i++ {
		clients[i] = <-provider.NewClients
	}

	_ = clients

	close(provider.EventStream)
}

func TestSendStringEventsNoReceiver(t *testing.T) {
	provider := NewProvider()

	testMsg1 := "ohy1Eo0i"
	testMsg2 := "fee3ahTh"
	testMsg3 := "ohb4ChuX"

	provider.EventStream <- testMsg1
	provider.EventStream <- testMsg2
	provider.EventStream <- testMsg3

	close(provider.EventStream)
}

func TestSendStringEvents(t *testing.T) {
	provider := NewProvider()
	const clients_num = 50
	var clients [clients_num]Client

	for i := 0; i < clients_num; i++ {
		clients[i] = <-provider.NewClients
	}

	// time for all the clients to be registrated properly
	<-time.After(50 * time.Millisecond)

	testMsg1 := "ohy1Eo0i"
	testMsg2 := "fee3ahTh"
	testMsg3 := "ohb4ChuX"

	provider.EventStream <- testMsg1
	for i := 0; i < clients_num; i++ {
		e := <-clients[i].Channel
		if s, ok := e.(string); ok {
			if s != testMsg1 {
				t.Error("Received wrong string")
			}
		} else {
			t.Error("Received no string, when one was expected.")
		}
	}

	provider.EventStream <- testMsg2
	for i := 0; i < clients_num; i++ {
		e := <-clients[i].Channel
		if s, ok := e.(string); ok {
			if s != testMsg2 {
				t.Error("Received wrong string")
			}
		} else {
			t.Error("Received no string, when one was expected.")
		}
	}

	provider.EventStream <- testMsg3
	for i := 0; i < clients_num; i++ {
		e := <-clients[i].Channel
		if s, ok := e.(string); ok {
			if s != testMsg3 {
				t.Error("Received wrong string")
			}
		} else {
			t.Error("Received no string, when one was expected.")
		}
	}

	close(provider.EventStream)
}

func TestSendComplexEvents(t *testing.T) {
	provider := NewProvider()
	const clients_num = 50
	var clients [clients_num]Client

	for i := 0; i < clients_num; i++ {
		clients[i] = <-provider.NewClients
	}

	// time for all the clients to be registrated properly
	<-time.After(50 * time.Millisecond)

	testMsg1 := struct{}{}
	testMsg2 := struct{ i int }{355}
	testMsg3 := struct{ k struct{} }{struct{}{}}

	provider.EventStream <- testMsg1
	for i := 0; i < clients_num; i++ {
		e := <-clients[i].Channel
		if s, ok := e.(struct{}); ok {
			if s != testMsg1 {
				t.Error("Received wrong Message")
			}
		} else {
			t.Error("Received no struct{}, when one was expected.")
		}
	}

	provider.EventStream <- testMsg2
	for i := 0; i < clients_num; i++ {
		e := <-clients[i].Channel
		if s, ok := e.(struct{ i int }); ok {
			if s != testMsg2 {
				t.Error("Received wrong message")
			}
		} else {
			t.Error("Received no struct{ i int }, when one was expected.")
		}
	}

	provider.EventStream <- testMsg3
	for i := 0; i < clients_num; i++ {
		e := <-clients[i].Channel
		if s, ok := e.(struct{ k struct{} }); ok {
			if s != testMsg3 {
				t.Error("Received wrong message")
			}
		} else {
			t.Error("Received no struct{ k struct{} }, when one was expected.")
		}
	}

	close(provider.EventStream)
}

func TestClientStrike(t *testing.T) {
	provider := NewProvider()
	referenceClient := <-provider.NewClients
	testClient := <-provider.NewClients

	// time for all the clients to be registrated properly
	<-time.After(50 * time.Millisecond)

	testMsg := "ohy1Eo0i"

	for i := 0; i < 6; i++ {
		provider.EventStream <- testMsg
		<-referenceClient.Channel
	}

	_, channelOpen := <-testClient.Channel
	if channelOpen {
		t.Error("Channel still open after 5 strikes")
	}

	close(provider.EventStream)
}

func TestClientQuit(t *testing.T) {
	provider := NewProvider()
	testClient := <-provider.NewClients

	// time for all the clients to be registrated properly
	<-time.After(50 * time.Millisecond)

	testMsg := "ohy1Eo0i"

	provider.EventStream <- testMsg
	e := <-testClient.Channel
	if s, ok := e.(string); ok {
		if s != testMsg {
			t.Error("Received wrong Message")
		}
	} else {
		t.Error("Received no struct{}, when one was expected.")
	}

	testClient.Quit()

	// time for the client to be unregistrated properly
	<-time.After(50 * time.Millisecond)

	provider.EventStream <- testMsg
	_, channelOpen := <-testClient.Channel

	if channelOpen {
		t.Error("Receive from channel ought to be closed")
	}

	close(provider.EventStream)
}

func TestComplexTest(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	provider := NewProvider()
	quit := make(chan struct{})

	go func() {
		for i := 0; ; i++ {
			select {
			case provider.EventStream <- i:
			case <-quit:
				return
			}
			<-time.After(10 * time.Millisecond)
		}
	}()

	for i := 0; i < 200; i++ {
		go func() {
			client := <-provider.NewClients
			for k := 0; k < 50; k++ {
				select {
				case z := <-client.Channel:
					if z == 2*i {
						client.Quit()
						return
					}
				case <-quit:
					return
				}
			}
		}()

		<-time.After(10 * time.Millisecond)
	}

	<-time.After(10 * time.Second)
	close(quit)
	<-time.After(50 * time.Millisecond)

	close(provider.EventStream)
}
