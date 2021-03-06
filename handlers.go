package main

import (
	"log"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type User struct {
	key string
}

type Room struct {
	name string
	host User
}

type UserForm struct {
	Name  string `form:"username" binding:"required"`
	Email string `form:"email"`
}

type IndexPage struct {
	Error string
}

type PingUserPage struct {
	Room   string
	IsHost bool
}

func indexGET(c *gin.Context) {
	c.Request.ParseForm()

	page := IndexPage{}
	errKey := c.Request.Form.Get("err")

	page.Error = errKey

	c.HTML(http.StatusOK, "index.tmpl", page)
}

func indexPOST(c *gin.Context) {
	var user UserForm
	c.BindWith(&user, binding.MultipartForm)

	key := c.MustGet("key").(string)

	name := url.QueryEscape(user.Name)

	// Check for room existance
	if exists("room:" + name) {
		log.Printf("Room exists")
		c.Redirect(http.StatusMovedPermanently, "/?err=taken")
		return
	}
	// Map new room name (user name) to user key (uuid)
	exec("SET", "room:"+name, key)
	exec("EXPIRE", "room:"+name, sessionExpire)
	// Map user to his room
	exec("SET", "user:"+key, name)
	exec("EXPIRE", "user:"+key, sessionExpire)

	c.Redirect(http.StatusMovedPermanently, "/ping/"+name)
}

func pingUserGET(c *gin.Context) {
	userKey := c.MustGet("key").(string)
	room := c.Params.ByName("user")

	roomUser, err := getString("GET", "room:"+room)
	if err != nil {
		log.Printf("Room %s doesn't exist\n", room)
		// If room doesn't exist, go to index
		c.Redirect(http.StatusMovedPermanently, "/")
		return
	}

	log.Printf("Room %s, host %s, user %s\n", room, roomUser, userKey)

	page := PingUserPage{room, false}

	c.Request.ParseForm()
	referer := c.Request.Form.Get("from")

	if roomUser != userKey {
		notifyHost(room, "notify:"+referer)
	} else {
		page.IsHost = true
	}

	log.Printf("%+v", page)

	c.HTML(http.StatusOK, "ping.tmpl", page)
}
