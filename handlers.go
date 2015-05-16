package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type User struct {
	Name  string `form:"username" binding:"required"`
	Email string `form:"email"`
}

func indexGET(c *gin.Context) {
	c.HTML(http.StatusOK, "index.tmpl", nil)
}

func indexPOST(c *gin.Context) {
	var user User
	c.BindWith(&user, binding.MultipartForm)

	key := c.MustGet("key").(string)

	// TODO: check for existance
	// Init new room
	//   - Save room name (user name) with user key (uuid)
	exec("SET", "room:"+user.Name, key)
	exec("SET", "user:"+key, user.Name)

	c.Redirect(http.StatusMovedPermanently, "/ping/"+user.Name)
}

func pingUserGET(c *gin.Context) {
	userKey := c.MustGet("key").(string)
	room := c.Params.ByName("user")

	roomUser, err := getString("GET", "room:"+room)
	if err != nil {
		// If room doesn't exist, go to index
		c.Redirect(http.StatusMovedPermanently, "/")
		return
	}

	log.Printf("Room %s, host %s, user %s\n", room, roomUser, userKey)
	if roomUser != "user:"+userKey {
		notifyHost(room)
	}

	c.HTML(http.StatusOK, "ping.tmpl", nil)
}
