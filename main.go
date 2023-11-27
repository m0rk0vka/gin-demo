package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func IndexHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "hello world",
	})
}

func main() {
	router := gin.Default()
	router.GET("/", IndexHandler)
	router.Run()
}
