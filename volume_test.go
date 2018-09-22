package lgtv

import (
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func TestGetVolume(t *testing.T) {
	server, url := startClientServer(t, func(c *websocket.Conn, req string) {
		if strings.Contains(req, "register") {
			c.WriteMessage(websocket.TextMessage, []byte("{\"id\":1,\"type\":\"registered\",\"payload\":{\"client-key\":\"123\"}}"))
		}
		if strings.Contains(req, uriVolumeGet) {
			c.WriteMessage(websocket.TextMessage, []byte("{\"id\":2,\"type\":\"response\",\"payload\":{\"volume\":10}}"))
		}
	})
	defer server.Close()

	lgtv := NewDefaultClient(url, key)
	err := lgtv.Connect()
	if err != nil {
		t.Error(err)
	}
	defer lgtv.Disconnect()

	v, err := lgtv.GetVolume()
	if err != nil {
		t.Error(err)
	}
	if v.Volume != 10 {
		t.Fail()
	}
}

func TestSetVolume(t *testing.T) {
	server, url := startClientServer(t, func(c *websocket.Conn, req string) {
		if strings.Contains(req, "register") {
			c.WriteMessage(websocket.TextMessage, []byte("{\"id\":1,\"type\":\"registered\",\"payload\":{\"client-key\":\"123\"}}"))
		}
		if strings.Contains(req, uriVolumeSet) {
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

	err = lgtv.SetVolume(20)
	if err != nil {
		t.Error(err)
	}
}

func TestIsMute(t *testing.T) {
	server, url := startClientServer(t, func(c *websocket.Conn, req string) {
		if strings.Contains(req, "register") {
			c.WriteMessage(websocket.TextMessage, []byte("{\"id\":1,\"type\":\"registered\",\"payload\":{\"client-key\":\"123\"}}"))
		}
		if strings.Contains(req, uriMuteGet) {
			c.WriteMessage(websocket.TextMessage, []byte("{\"id\":2,\"type\":\"response\",\"payload\":{\"mute\":true}}"))
		}
	})
	defer server.Close()

	lgtv := NewDefaultClient(url, key)
	err := lgtv.Connect()
	if err != nil {
		t.Error(err)
	}
	defer lgtv.Disconnect()

	m, err := lgtv.IsMute()
	if err != nil {
		t.Error(err)
	}
	if !m {
		t.Fail()
	}
}

func TestSetMute(t *testing.T) {
	server, url := startClientServer(t, func(c *websocket.Conn, req string) {
		if strings.Contains(req, "register") {
			c.WriteMessage(websocket.TextMessage, []byte("{\"id\":1,\"type\":\"registered\",\"payload\":{\"client-key\":\"123\"}}"))
		}
		if strings.Contains(req, uriMuteSet) {
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

	err = lgtv.SetMute(false)
	if err != nil {
		t.Error(err)
	}
}

func TestSubscribeVolume(t *testing.T) {
	server, url := startClientServer(t, func(c *websocket.Conn, req string) {
		if strings.Contains(req, "register") {
			c.WriteMessage(websocket.TextMessage, []byte("{\"id\":1,\"type\":\"registered\",\"payload\":{\"client-key\":\"123\"}}"))
		}
		if strings.Contains(req, "subscribe") && strings.Contains(req, uriVolumeGet) {
			go func() {
				time.Sleep(20 * time.Millisecond)
				c.WriteMessage(websocket.TextMessage, []byte("{\"id\":2,\"type\":\"response\",\"payload\":{\"volume\":10}}"))
				time.Sleep(50 * time.Millisecond)
				c.WriteMessage(websocket.TextMessage, []byte("{\"id\":2,\"type\":\"response\",\"payload\":{\"volume\":20}}"))
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

	ch := make(chan int)
	err = lgtv.SubscribeVolume(func(v Volume) {
		ch <- v.Volume
	})
	if err != nil {
		t.Error(err)
	}

	for i := 0; i < 2; i++ {
		select {
		case volume := <-ch:
			if i == 0 && volume != 10 {
				t.Logf("First update should be 10, but is %d", volume)
				t.Fail()
			}
			if i == 1 && volume != 20 {
				t.Logf("Second update should be 20, but is %d", volume)
				t.Fail()
			}
		case <-time.After(100 * time.Millisecond):
			t.Log("Timeout")
			t.Fail()
		}
	}
}
