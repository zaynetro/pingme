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

	log.Printf("user name %s\n", user.Name)

	c.Redirect(http.StatusMovedPermanently, "/ping/"+user.Name)
}

func pingUserGET(c *gin.Context) {
	c.HTML(http.StatusOK, "ping.tmpl", nil)
	// c.String(http.StatusOK, "pinged "+c.Params.ByName("user"))
}
