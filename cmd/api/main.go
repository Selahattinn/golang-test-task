package main

import (
	"github.com/gin-gonic/gin"
	rabbit "twitch_chat_analysis/pkg/rabbitmq"
)

func main() {
	r := gin.Default()

	err := rabbit.ConnRabbitMQ()
	if err != nil {
		panic(err)
	}

	r.GET("/test", func(c *gin.Context) {
		c.JSON(200, "worked")
	})

	r.POST("/message", func(c *gin.Context) {
		var message rabbit.Message
		err := c.ShouldBindJSON(&message)
		if err != nil {
			c.JSON(400, err.Error())
			return
		}

		err = rabbit.Publish(message)
		if err != nil {
			c.JSON(400, err.Error())
			return
		} else {
			c.JSON(200, "")
			return

		}

	})

	r.Run()
}
