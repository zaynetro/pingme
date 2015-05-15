package main

import (
	"log"
	"net/http"
	"time"

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

var clients map[string]*websocket.Conn

func wshandler(w http.ResponseWriter, r *http.Request) {
	clients = make(map[string]*websocket.Conn)

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to set websocket upgrade: %+v\n", err)
		return
	}

	clients["roman"] = conn

	for {
		t, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}
		log.Printf("Socket message type: %d", t)

		conn.WriteMessage(t, msg)
	}
}
