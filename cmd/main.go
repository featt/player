package main

import (
	"log"
	"net/http"

	c "github.com/featt/player/pkg/controllers"
	"github.com/gin-gonic/gin"
)


func main() {

	router := gin.Default()

    router.POST("/add-song", c.AddSong)
	router.GET("/play", c.Play)
	router.GET("/pause", c.Pause)
	router.GET("/prev", c.Prev)
	router.GET("/next", c.Next)

	log.Fatal(http.ListenAndServe(":8080", router))
}