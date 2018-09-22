package lgtv

import (
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func TestGetCurrentChannel(t *testing.T) {
	server, url := startClientServer(t, func(c *websocket.Conn, req string) {
		if strings.Contains(req, "register") {
			c.WriteMessage(websocket.TextMessage, []byte("{\"id\":1,\"type\":\"registered\",\"payload\":{\"client-key\":\"123\"}}"))
		}
		if strings.Contains(req, uriChannelGet) {
			c.WriteMessage(websocket.TextMessage, []byte("{\"id\":2,\"type\":\"response\",\"payload\":{\"channelID\":\"foo\"}}"))
		}
	})
	defer server.Close()

	lgtv := NewDefaultClient(url, key)
	err := lgtv.Connect()
	if err != nil {
		t.Error(err)
	}
	defer lgtv.Disconnect()

	channel, err := lgtv.GetCurrentChannel()
	if err != nil {
		t.Error(err)
	}
	if channel.ID != "foo" {
		t.Fail()
	}
}

func TestGetChannels(t *testing.T) {
	server, url := startClientServer(t, func(c *websocket.Conn, req string) {
		if strings.Contains(req, "register") {
			c.WriteMessage(websocket.TextMessage, []byte("{\"id\":1,\"type\":\"registered\",\"payload\":{\"client-key\":\"123\"}}"))
		}
		if strings.Contains(req, uriChannelList) {
			c.WriteMessage(websocket.TextMessage, []byte("{\"id\":2,\"type\":\"response\",\"payload\":{\"channels\":[{\"channelID\":\"foo\"},{\"channelID\":\"bar\"}]}}"))
		}
	})
	defer server.Close()

	lgtv := NewDefaultClient(url, key)
	err := lgtv.Connect()
	if err != nil {
		t.Error(err)
	}
	defer lgtv.Disconnect()

	list, err := lgtv.GetChannels()
	if err != nil {
		t.Error(err)
	}
	if len(list) != 2 {
		t.Fail()
		return
	}
	if list[0].ID != "foo" {
		t.Fail()
	}
	if list[1].ID != "bar" {
		t.Fail()
	}
}

func TestSetChannel(t *testing.T) {
	server, url := startClientServer(t, func(c *websocket.Conn, req string) {
		if strings.Contains(req, "register") {
			c.WriteMessage(websocket.TextMessage, []byte("{\"id\":1,\"type\":\"registered\",\"payload\":{\"client-key\":\"123\"}}"))
		}
		if strings.Contains(req, uriChannelSet) {
			c.WriteMessage(websocket.TextMessage, []byte("{\"id\":2,\"type\":\"response\",\"payload\":null}"))
		}
	})
	defer server.Close()

	lgtv := NewDefaultClient(url, key)
	err := lgtv.Connect()
	if err != nil {
		t.Error(err)
	}
	defer lgtv.Disconnect()

	err = lgtv.SetChannel("42")
	if err != nil {
		t.Error(err)
	}
}

func TestSubscribeChannel(t *testing.T) {
	server, url := startClientServer(t, func(c *websocket.Conn, req string) {
		if strings.Contains(req, "register") {
			c.WriteMessage(websocket.TextMessage, []byte("{\"id\":1,\"type\":\"registered\",\"payload\":{\"client-key\":\"123\"}}"))
		}
		if strings.Contains(req, "subscribe") && strings.Contains(req, uriChannelGet) {
			go func() {
				time.Sleep(20 * time.Millisecond)
				c.WriteMessage(websocket.TextMessage, []byte("{\"id\":2,\"type\":\"response\",\"payload\":{\"channelID\":\"foo\"}}"))
				time.Sleep(50 * time.Millisecond)
				c.WriteMessage(websocket.TextMessage, []byte("{\"id\":2,\"type\":\"response\",\"payload\":{\"channelID\":\"bar\"}}"))
			}()
		}
	})
	defer server.Close()

	lgtv := NewDefaultClient(url, key)
	err := lgtv.Connect()
	if err != nil {
		t.Error(err)
	}
	defer lgtv.Disconnect()

	ch := make(chan string)
	err = lgtv.SubscribeChannel(func(c Channel) {
		ch <- c.ID
	})
	if err != nil {
		t.Error(err)
	}

	for i := 0; i < 2; i++ {
		select {
		case channelID := <-ch:
			if i == 0 && channelID != "foo" {
				t.Logf("First update should be \"foo\", but is \"%s\"", channelID)
				t.Fail()
			}
			if i == 1 && channelID != "bar" {
				t.Logf("Second update should be \"bar\", but is \"%s\"", channelID)
				t.Fail()
			}
		case <-time.After(100 * time.Millisecond):
			t.Log("Timeout")
			t.Fail()
		}
	}
}
