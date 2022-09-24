package main

import (
	"context"
	rabbit "twitch_chat_analysis/pkg/rabbitmq"
	"twitch_chat_analysis/pkg/redis"
)

func main() {
	err := rabbit.ConnRabbitMQ()
	if err != nil {
		panic(err)
	}

	r, err := redis.InitRedis(redis.Config{
		URL:      "redis",
		Password: "",
		Database: 0,
		Port:     "6379",
		TLS:      false,
	})

	if err != nil {
		panic(err)
	}

	ch, err := rabbit.Consume()
	if err != nil {
		panic(err)

	}

	for msg := range ch {
		ctx := context.Background()
		r.Set(ctx, msg.MessageId, msg.Body, 0)
	}
}
