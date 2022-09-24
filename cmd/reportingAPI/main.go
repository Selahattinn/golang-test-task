package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"twitch_chat_analysis/pkg/redis"
)

func main() {
	r := gin.Default()

	re, err := redis.InitRedis(redis.Config{
		URL:      "redis",
		Password: "",
		Database: 0,
		Port:     "6379",
		TLS:      false,
	})

	if err != nil {
		panic(err)
	}

	r.GET("/message/list", func(c *gin.Context) {
		ctx := context.Background()
		//get params
		sender := c.Query("sender")
		if sender == "" {
			c.JSON(400, "sender is required")
			return
		}

		receiver := c.Query("receiver")
		if receiver == "" {
			c.JSON(400, "receiver is required")
			return
		}

		messages, err := re.Get(ctx, sender+receiver+"*")
		if err != nil {
			c.JSON(400, err.Error())
			return
		}

		c.JSON(200, messages)
		return

	})

	r.Run()
}
