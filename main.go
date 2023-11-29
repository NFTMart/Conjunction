package main

import (
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

func main() {
	LoadNodes("nodes.json")
	r := gin.Default()
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	r.Use(cors.New(config))
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(gzip.Gzip(gzip.DefaultCompression))

	Router(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8129"
		log.Printf("defaulting to port %s", port)
	}
	err := r.Run(":" + port)
	if err != nil {
		log.Print(err)
	}
}
