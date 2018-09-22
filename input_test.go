package lgtv

import (
	"strings"
	"testing"

	"github.com/gorilla/websocket"
)

func TestGetInputs(t *testing.T) {
	server, url := startClientServer(t, func(c *websocket.Conn, req string) {
		if strings.Contains(req, "register") {
			c.WriteMessage(websocket.TextMessage, []byte("{\"id\":1,\"type\":\"registered\",\"payload\":{\"client-key\":\"123\"}}"))
		}
		if strings.Contains(req, uriInputList) {
			c.WriteMessage(websocket.TextMessage, []byte("{\"id\":2,\"type\":\"response\",\"payload\":{\"devices\":[{\"id\":\"foo\"},{\"id\":\"bar\"}]}}"))
		}
	})
	defer server.Close()

	lgtv := NewDefaultClient(url, key)
	err := lgtv.Connect()
	if err != nil {
		t.Error(err)
	}
	defer lgtv.Disconnect()

	list, err := lgtv.GetInputs()
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

func TestSwitchInput(t *testing.T) {
	server, url := startClientServer(t, func(c *websocket.Conn, req string) {
		if strings.Contains(req, "register") {
			c.WriteMessage(websocket.TextMessage, []byte("{\"id\":1,\"type\":\"registered\",\"payload\":{\"client-key\":\"123\"}}"))
		}
		if strings.Contains(req, uriInputSet) {
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

	err = lgtv.SwitchInput("foo")
	if err != nil {
		t.Error(err)
	}
}
