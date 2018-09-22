package connection

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
)

func (c *Connection) ping() {
	ticker := time.NewTicker(c.timeouts.Ping)
	defer ticker.Stop()

	for c.IsConnected() {
		select {
		case <-ticker.C:
			c.mutex.Lock()
			if c.IsConnected() {
				log.Println("-> ping")
				if err := c.ws.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(c.timeouts.Write)); err != nil {
					log.Printf("ping write: %s", err)
					c.mutex.Unlock()
					c.disconnect(err)
					break
				}
			}
			c.mutex.Unlock()
		case <-c.done:
			break
		}
	}
}

func (c *Connection) pong(string) error {
	log.Println("<- pong")
	return c.ws.SetReadDeadline(time.Now().Add(c.timeouts.Read))
}
