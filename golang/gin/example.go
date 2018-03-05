package main

import "github.com/gin-gonic/gin"

func run1() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run(":8001")
}

func run2() {
	router := gin.Default()

	v1 := router.Group("/v1")
	v1.POST("/form_post", func(c *gin.Context) {
		message := c.PostForm("message")
		c.JSON(200, gin.H{
			"status":  "posted",
			"message": message,
		})
	})

	router.Run(":8080")
}

func main() {
	run2()
}
