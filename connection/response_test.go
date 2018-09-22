package connection

import (
	"net/http"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func TestRequestWithErrorResponse(t *testing.T) {
	server, url := startServer(func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer c.Close()
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				break
			}
			if string(message) == "{\"id\":1,\"type\":\"request\",\"uri\":\"/foo\",\"payload\":{\"prop1\":\"test\"}}\n" {
				c.WriteMessage(websocket.TextMessage, []byte("{\"id\":1,\"type\":\"error\",\"error\":\"foo\",\"payload\":null}"))
			}
		}
	})
	defer server.Close()

	conn := NewConnection(url, DefaultTimeouts())

	err := conn.Connect()
	if err != nil {
		t.Error(err)
	}
	defer conn.Close()

	err = conn.Request(1, "/foo", testPayload{"test"}, nil)
	if err == nil {
		t.Fail()
	}
}

func TestRequestWithInvalidResponsePayload(t *testing.T) {
	server, url := startServer(func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer c.Close()
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				break
			}
			if string(message) == "{\"id\":1,\"type\":\"request\",\"uri\":\"/foo\",\"payload\":{\"prop1\":\"test\"}}\n" {
				c.WriteMessage(websocket.TextMessage, []byte("{\"id\":1,\"type\":\"response\",\"error\":null,\"payload\":{\"prop1\":1}}"))
			}
		}
	})
	defer server.Close()

	conn := NewConnection(url, DefaultTimeouts())

	err := conn.Connect()
	if err != nil {
		t.Error(err)
	}
	defer conn.Close()

	result := testPayload{}
	err = conn.Request(1, "/foo", testPayload{"test"}, &result)
	if err == nil {
		t.Fail()
	}
}

func TestRequestWithResponseTimeout(t *testing.T) {
	server, url := startServer(func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer c.Close()
		for {
			_, _, err := c.ReadMessage()
			if err != nil {
				break
			}
			time.Sleep(20 * time.Millisecond)
		}
	})
	defer server.Close()

	timeouts := Timeouts{
		Read:  100 * time.Millisecond,
		Ping:  50 * time.Millisecond,
		Write: 100 * time.Millisecond,
	}
	conn := NewConnection(url, timeouts)

	err := conn.Connect()
	if err != nil {
		t.Error(err)
	}
	defer conn.Close()

	err = conn.Request(1, "/foo", testPayload{"test"}, nil)
	if err == nil {
		t.Fail()
	}
	time.Sleep(500 * time.Millisecond)
}
