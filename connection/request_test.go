package connection

import (
	"net/http"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

type testPayload struct {
	Prop1 string `json:"prop1"`
}

func TestRequest(t *testing.T) {
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
				c.WriteMessage(websocket.TextMessage, []byte("{\"id\":1,\"type\":\"response\",\"error\":null,\"payload\":{\"prop1\":\"foo\"}}"))
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
	if err != nil {
		t.Error(err)
	}
	if result.Prop1 != "foo" {
		t.Fail()
	}
}

func TestRequestWithoutResponsePayload(t *testing.T) {
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
				c.WriteMessage(websocket.TextMessage, []byte("{\"id\":1,\"type\":\"response\",\"error\":null,\"payload\":{\"prop1\":\"foo\"}}"))
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
	if err != nil {
		t.Error(err)
	}
}

func TestSubscribe(t *testing.T) {
	server, url := startServer(func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer c.Close()
		for {
			c.WriteMessage(websocket.TextMessage, []byte("{\"id\":1,\"type\":\"response\",\"error\":null,\"payload\":{\"prop1\":\"foo\"}}"))
			_, _, err := c.ReadMessage()
			if err != nil {
				break
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

	ch, err := conn.Subscribe(1, "/foo")
	if err != nil {
		t.Error(err)
	}

	select {
	case <-time.After(100 * time.Millisecond):
		t.Log("Timeout")
		t.Fail()
	case resp := <-ch:
		if resp.Type != "response" {
			t.Fail()
		}
	}
}
