package lgtv

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
)

const key = "123"

var upgrader = websocket.Upgrader{}

func startClientServer(t *testing.T, handler func(*websocket.Conn, string)) (*httptest.Server, string) {
	f := func(w http.ResponseWriter, r *http.Request) {
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
			handler(c, string(message))
		}
	}

	s := httptest.NewServer(http.HandlerFunc(f))

	return s, "ws" + strings.TrimPrefix(s.URL, "http")
}

func TestConnect(t *testing.T) {
	server, url := startClientServer(t, func(c *websocket.Conn, req string) {
		if strings.Contains(req, "register") {
			c.WriteMessage(websocket.TextMessage, []byte("{\"id\":1,\"type\":\"registered\",\"payload\":{\"client-key\":\"123\"}}"))
		}
	})
	defer server.Close()

	connected := false
	disconnected := false

	lgtv := NewDefaultClient(url, key)
	lgtv.OnConnect(func(*Client) {
		connected = true
	})
	lgtv.OnDisconnect(func(*Client, error) {
		disconnected = true
	})

	err := lgtv.Connect()
	if err != nil {
		t.Error(err)
	}
	defer lgtv.Disconnect()

	if !connected {
		t.Log("connected should be true, but is false")
		t.Fail()
	}
	if !lgtv.IsConnected() {
		t.Log("IsConnected() returns false, but should return true")
		t.Fail()
	}
	if disconnected {
		t.Log("OnDisconnect was called")
		t.Fail()
	}
}

func TestConnectWithFailedRegistration(t *testing.T) {
	server, url := startClientServer(t, func(c *websocket.Conn, req string) {
		if strings.Contains(req, "register") {
			c.WriteMessage(websocket.TextMessage, []byte("{\"id\":1,\"type\":\"error\"}"))
		}
	})
	defer server.Close()

	connected := false
	disconnected := false

	lgtv := NewDefaultClient(url, key)
	lgtv.OnConnect(func(*Client) {
		connected = true
	})
	lgtv.OnDisconnect(func(*Client, error) {
		disconnected = true
	})

	err := lgtv.Connect()
	if err == nil {
		defer lgtv.Disconnect()
		t.Fail()
	}

	if connected {
		t.Log("connected should be false, but is true")
		t.Fail()
	}
	if lgtv.IsConnected() {
		t.Log("IsConnected() returns true, but should return false")
		t.Fail()
	}
	if !disconnected {
		t.Log("OnDisconnect was not called")
		t.Fail()
	}
}

func TestRequestWithoutConnection(t *testing.T) {
	server, url := startClientServer(t, func(c *websocket.Conn, req string) {})
	defer server.Close()

	lgtv := NewDefaultClient(url, key)
	err := lgtv.SetMute(true)
	if err == nil {
		t.Fail()
	}
}
