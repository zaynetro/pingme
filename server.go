package main

import (
	"log"
	"os"
	"runtime"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/twinj/uuid"
)

func main() {
	nuCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(nuCPU)

	log.SetOutput(os.Stdout)
	log.Printf("Running with %d CPUs\n", nuCPU)
	log.Printf("CONFIG=%v\n", CONFIG)
	// Configure web server
	SetUpServer()
}

func SetUpServer() {
	r := gin.Default()

	store, _ := sessions.NewRedisStore(10, "tcp", CONFIG.Redis.Host, CONFIG.Redis.Password, []byte(CONFIG.Server.Secret))
	r.Use(sessions.Sessions("session", store))

	r.Use(gin.Recovery())
	r.Use(Uniquify())

	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "public/")

	SetUpRoutes(r)

	r.GET("/incr", func(c *gin.Context) {
		c.JSON(200, gin.H{"count": c.MustGet("key").(string)})
	})

	// var port string

	r.Run(":" + CONFIG.Server.Port)
}

func SetUpRoutes(r *gin.Engine) {
	r.GET("/", indexGET)
	r.POST("/", indexPOST)
	r.GET("/ping/:user", pingUserGET)

	// WebSockets
	r.GET("/ws", wshandler)
}

// Create uuid key for each user session
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
