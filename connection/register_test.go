package connection

import (
	"net/http"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
)

func TestRegister(t *testing.T) {
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
			if strings.Contains(string(message), "register") {
				c.WriteMessage(websocket.TextMessage, []byte("{\"id\":1,\"type\":\"registered\",\"payload\":{\"client-key\":\"42\"}}"))
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

	err = conn.Register(1, "42", make([]string, 0))
	if err != nil {
		t.Error(err)
	}
}

func TestRegisterWithoutKey(t *testing.T) {
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
			if strings.Contains(string(message), "register") {
				c.WriteMessage(websocket.TextMessage, []byte("{\"id\":1,\"type\":\"response\"}"))
				c.WriteMessage(websocket.TextMessage, []byte("{\"id\":1,\"type\":\"registered\",\"payload\":{\"client-key\":\"42\"}}"))
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

	err = conn.Register(1, "", make([]string, 0))
	if err != nil {
		t.Error(err)
	}
}
