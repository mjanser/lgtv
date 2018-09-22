package lgtv

import (
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func TestGetApps(t *testing.T) {
	server, url := startClientServer(t, func(c *websocket.Conn, req string) {
		if strings.Contains(req, "register") {
			c.WriteMessage(websocket.TextMessage, []byte("{\"id\":1,\"type\":\"registered\",\"payload\":{\"client-key\":\"123\"}}"))
		}
		if strings.Contains(req, uriAppList) {
			c.WriteMessage(websocket.TextMessage, []byte("{\"id\":2,\"type\":\"response\",\"payload\":{\"apps\":[{\"id\":\"foo\"},{\"id\":\"bar\"}]}}"))
		}
	})
	defer server.Close()

	lgtv := NewDefaultClient(url, key)
	err := lgtv.Connect()
	if err != nil {
		t.Error(err)
	}
	defer lgtv.Disconnect()

	list, err := lgtv.GetApps()
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

func TestLaunchApp(t *testing.T) {
	server, url := startClientServer(t, func(c *websocket.Conn, req string) {
		if strings.Contains(req, "register") {
			c.WriteMessage(websocket.TextMessage, []byte("{\"id\":1,\"type\":\"registered\",\"payload\":{\"client-key\":\"123\"}}"))
		}
		if strings.Contains(req, uriAppLaunch) {
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

	err = lgtv.LaunchApp("foo")
	if err != nil {
		t.Error(err)
	}
}

func TestSubscribeApp(t *testing.T) {
	server, url := startClientServer(t, func(c *websocket.Conn, req string) {
		if strings.Contains(req, "register") {
			c.WriteMessage(websocket.TextMessage, []byte("{\"id\":1,\"type\":\"registered\",\"payload\":{\"client-key\":\"123\"}}"))
		}
		if strings.Contains(req, "subscribe") && strings.Contains(req, uriAppGet) {
			go func() {
				time.Sleep(20 * time.Millisecond)
				c.WriteMessage(websocket.TextMessage, []byte("{\"id\":2,\"type\":\"response\",\"payload\":{\"appId\":\"foo\"}}"))
				time.Sleep(50 * time.Millisecond)
				c.WriteMessage(websocket.TextMessage, []byte("{\"id\":2,\"type\":\"response\",\"payload\":{\"appId\":\"bar\"}}"))
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
	err = lgtv.SubscribeApp(func(a AppInfo) {
		ch <- a.ID
	})
	if err != nil {
		t.Error(err)
	}

	for i := 0; i < 2; i++ {
		select {
		case appID := <-ch:
			if i == 0 && appID != "foo" {
				t.Logf("First update should be \"foo\", but is \"%s\"", appID)
				t.Fail()
			}
			if i == 1 && appID != "bar" {
				t.Logf("Second update should be \"bar\", but is \"%s\"", appID)
				t.Fail()
			}
		case <-time.After(100 * time.Millisecond):
			t.Log("Timeout")
			t.Fail()
		}
	}
}
