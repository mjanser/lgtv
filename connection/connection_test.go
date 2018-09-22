package connection

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}
var sendMessageChan = make(chan []byte)

func startServer(f func(http.ResponseWriter, *http.Request)) (*httptest.Server, string) {
	if f == nil {
		f = func(w http.ResponseWriter, r *http.Request) {
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
			}
		}
	}
	s := httptest.NewServer(http.HandlerFunc(f))

	return s, "ws" + strings.TrimPrefix(s.URL, "http")
}

func TestConnect(t *testing.T) {
	server, url := startServer(nil)
	defer server.Close()

	connected := false
	disconnected := false
	timeouts := Timeouts{
		Read: 20 * time.Millisecond,
		Ping: 10 * time.Millisecond,
	}

	conn := NewConnection(url, timeouts)
	conn.OnConnect = func(c *Connection) {
		connected = true
	}
	conn.OnDisconnect = func(c *Connection, err error) {
		disconnected = true
	}

	err := conn.Connect()
	if err != nil {
		t.Error(err)
	}
	defer conn.Close()

	if !connected {
		t.Log("connected should be true, but is false")
		t.Fail()
	}
	if !conn.IsConnected() {
		t.Log("IsConnected() returns false, but should return true")
		t.Fail()
	}
	if disconnected {
		t.Log("OnDisconnect was called")
		t.Fail()
	}
}

func TestConnectError(t *testing.T) {
	server, _ := startServer(nil)
	defer server.Close()

	conn := NewConnection("ws://localhost:12000", Timeouts{})

	err := conn.Connect()
	if err == nil {
		t.Fail()
	}
	defer conn.Close()
}

func TestServerGoneAway(t *testing.T) {
	server, url := startServer(nil)
	defer server.Close()

	connected := false
	disconnected := false
	timeouts := Timeouts{
		Read:  20 * time.Millisecond,
		Ping:  10 * time.Millisecond,
		Write: 20 * time.Millisecond,
	}

	conn := NewConnection(url, timeouts)
	conn.OnConnect = func(c *Connection) {
		connected = true
	}
	conn.OnDisconnect = func(c *Connection, err error) {
		disconnected = true
	}

	err := conn.Connect()
	if err != nil {
		t.Error(err)
	}
	defer conn.Close()

	time.Sleep(30 * time.Millisecond)
	conn.ws.Close()
	time.Sleep(30 * time.Millisecond)

	if !connected {
		t.Log("connected should be true, but is false")
		t.Fail()
	}
	if conn.IsConnected() {
		t.Log("IsConnected() returns true, but should return false")
		t.Fail()
	}
	if !disconnected {
		t.Log("OnDisconnect was not called")
		t.Fail()
	}
}
