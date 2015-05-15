package main

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type connection struct {
	ws *websocket.Conn
	// queue chan []byte
}

type hub struct {
	rooms map[string]*connection
}

var h = hub{
	rooms: make(map[string]*connection),
}

func wshandler(c *gin.Context) {
	w := c.Writer
	r := c.Request

	userKey := c.MustGet("key").(string)

	log.Printf("User key: %s", userKey)

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to set websocket upgrade: %+v\n", err)
		return
	}

	conLink := &connection{ws: conn}

	h.rooms[userKey] = conLink

	for {
		t, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}
		log.Printf("Socket message type: %d", t)

		conn.WriteMessage(t, msg)
	}
}
