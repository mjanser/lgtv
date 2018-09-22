package connection

import (
	"testing"
	"time"
)

func TestPingPong(t *testing.T) {
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

	time.Sleep(40 * time.Millisecond)

	if !connected {
		t.Log("connected should be false, but is true")
		t.Fail()
	}
	if !conn.IsConnected() {
		t.Log("isConnected() returns false, but should return true")
		t.Fail()
	}
	if disconnected {
		t.Log("onDisconnect was called")
		t.Fail()
	}
}

func TestReadTimeout(t *testing.T) {
	server, url := startServer(nil)
	defer server.Close()

	connected := false
	disconnected := false
	timeouts := Timeouts{
		Read:  20 * time.Millisecond,
		Ping:  40 * time.Millisecond,
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

	time.Sleep(40 * time.Millisecond)

	if !connected {
		t.Log("connected should be false, but is true")
		t.Fail()
	}
	if conn.IsConnected() {
		t.Log("isConnected() returns true, but should return false")
		t.Fail()
	}
	if !disconnected {
		t.Log("onDisconnect was not called")
		t.Fail()
	}
}
