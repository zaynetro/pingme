package main

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

const (
	pingPeriod = 60 * time.Second
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type connection struct {
	ws *websocket.Conn
	// send chan []byte
	user *User
	room *Room
}

type hub struct {
	rooms map[string]*connection
}

var h = hub{
	rooms: make(map[string]*connection),
}

func wshandler(context *gin.Context) {
	w := context.Writer
	r := context.Request

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to set websocket upgrade: %+v\n", err)
		return
	}

	var room *Room
	user := User{context.MustGet("key").(string)}

	userRoom, err := getString("GET", "user:"+user.key)
	// If user has a room save his connection
	if err == nil {
		room = &Room{userRoom, user}
		log.Printf("User key: %s\n", user.key)
		log.Printf("User room: %v\n", room)
	}

	c := &connection{
		ws: conn,
		// send: make(chan []byte),
		user: &user,
		room: room,
	}

	if err == nil {
		// Save connection only when user has a room
		h.rooms[room.name] = c
	}

	c.manage()
}

func (c *connection) manage() {
	// ticker := time.NewTicker(pingPeriod)

	log.Printf("INIT CONNECTION with USER=%v and ROOM=%v\n", c.user, c.room)

	defer func() {
		// Remove room and connection when host left
		if c.room != nil {
			delete(h.rooms, c.room.name)
			exec("DEL", "user:"+c.user.key, "room:"+c.room.name)
			log.Printf("Room removed")
		}
		// ticker.Stop()
		c.ws.Close()
		log.Printf("CLOSE CONNECTION\n")
	}()

	for {
		// select {
		// case <-ticker.C:
		// 	c.ws.WriteMessage(websocket.PingMessage, []byte("ping:empty"))
		// 	break
		// default:
		// Echo received message
		t, msg, err := c.ws.ReadMessage()
		if err != nil {
			break
		}

		c.ws.WriteMessage(t, msg)
		// }
	}

}

func notifyHost(host string, msg string) {
	log.Printf("Notify host\n")

	c, ok := h.rooms[host]
	if !ok {
		log.Printf("No connection for host\n")
		return
	}

	c.ws.WriteMessage(websocket.TextMessage, []byte(msg))
}
