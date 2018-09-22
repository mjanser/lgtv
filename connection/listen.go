package connection

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

func (c *Connection) listen() {
	for c.IsConnected() {
		_, message, err := c.ws.ReadMessage()
		if err != nil {
			log.Printf("read: %s", err)
			if !c.IsConnected() {
				break
			}
			if !websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				c.disconnect(err)
				break
			}
			c.disconnect(nil)
			break
		} else {
			resp := Response{}
			if err := json.Unmarshal(message, &resp); err != nil {
				log.Printf("Could not unmarshal response \"%s\", error was: %s", string(message), err)
				continue
			}

			log.Printf("<- %#v", string(message))

			if ch, ok := c.responseChans[resp.ID]; ok {
				ch <- resp
			} else if ch, ok := c.notificationChans[resp.ID]; ok {
				ch <- resp
			} else {
				log.Printf("Could not assign response with ID %d (%s) to any listener", resp.ID, string(message))
			}
		}
	}
}
