package main

import (
	"flag"
	"log"
	"os"
	"runtime"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/twinj/uuid"
)

var addr = flag.String("addr", ":3000", "server address")

func main() {
	flag.Parse()

	nuCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(nuCPU)

	log.SetOutput(os.Stdout)
	log.Printf("Running with %d CPUs\n", nuCPU)

	SetUpServer()
}

func SetUpServer() {
	r := gin.Default()

	store, _ := sessions.NewRedisStore(10, "tcp", "localhost:6379", "", []byte("secret"))
	r.Use(sessions.Sessions("session", store))

	r.Use(gin.Recovery())
	r.Use(Uniquify())

	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "public/")

	SetUpRoutes(r)

	r.GET("/incr", func(c *gin.Context) {
		c.JSON(200, gin.H{"count": c.MustGet("key").(string)})
	})

	r.Run(*addr)
}

func SetUpRoutes(r *gin.Engine) {
	r.GET("/", indexGET)
	r.POST("/", indexPOST)
	r.GET("/ping/:user", pingUserGET)

	// WebSockets
	r.GET("/ws", wshandler)
}

func Uniquify() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		key := session.Get("key")
		if key == nil {
			key = uuid.NewV4().String()
			session.Set("key", key)
			session.Save()
		}

		c.Set("key", key)
		c.Next()
	}
}
