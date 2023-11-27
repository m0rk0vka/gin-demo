package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func IndexHandler(c *gin.Context) {
	name := c.Params.ByName("name")
	c.JSON(http.StatusOK, gin.H{
		"message": "hello " + name,
	})
}

func main() {
	router := gin.Default()
	router.GET("/:name", IndexHandler)
	router.Run()
}
