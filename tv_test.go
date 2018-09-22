package lgtv

import (
	"strings"
	"testing"

	"github.com/gorilla/websocket"
)

func TestToast(t *testing.T) {
	server, url := startClientServer(t, func(c *websocket.Conn, req string) {
		if strings.Contains(req, "register") {
			c.WriteMessage(websocket.TextMessage, []byte("{\"id\":1,\"type\":\"registered\",\"payload\":{\"client-key\":\"123\"}}"))
		}
		if strings.Contains(req, uriToast) {
			c.WriteMessage(websocket.TextMessage, []byte("{\"id\":2,\"type\":\"response\",\"payload\":{}}"))
		}
	})
	defer server.Close()

	lgtv := NewDefaultClient(url, key)
	err := lgtv.Connect()
	if err != nil {
		t.Error(err)
	}
	defer lgtv.Disconnect()

	err = lgtv.Toast("foo")
	if err != nil {
		t.Error(err)
	}
}

func TestTurnOff(t *testing.T) {
	server, url := startClientServer(t, func(c *websocket.Conn, req string) {
		if strings.Contains(req, "register") {
			c.WriteMessage(websocket.TextMessage, []byte("{\"id\":1,\"type\":\"registered\",\"payload\":{\"client-key\":\"123\"}}"))
		}
		if strings.Contains(req, uriPowerOff) {
			c.WriteMessage(websocket.TextMessage, []byte("{\"id\":2,\"type\":\"response\",\"payload\":{}}"))
		}
	})
	defer server.Close()

	lgtv := NewDefaultClient(url, key)
	err := lgtv.Connect()
	if err != nil {
		t.Error(err)
	}
	defer lgtv.Disconnect()

	err = lgtv.TurnOff()
	if err != nil {
		t.Error(err)
	}
}
