package rabbit

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"twitch_chat_analysis/pkg/constants"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Message struct {
	Sender   string `json:"sender" binding:"required"`
	Receiver string `json:"receiver" binding:"required"`
	Message  string `json:"message" binding:"required"`
}

// var conn *amqp.Connection
var channel *amqp.Channel

func ConnRabbitMQ() error {
	conn, err := amqp.Dial("amqp://user:password@rabbitmq:5672/")
	if err != nil {
		return err
	}

	channel, err = conn.Channel()
	if err != nil {
		return err
	}

	channel.QueueDeclare(
		constants.RabbitChannelName, // name
		true,                        // durable
		false,                       // delete when unused
		false,                       // exclusive
		false,                       // no-wait
		nil,                         // arguments
	)

	return nil
}

func Publish(m Message) error {
	//check for fields empty
	if m.Message == "" || m.Sender == "" || m.Receiver == "" {
		return nil
	}

	err := channel.PublishWithContext(
		context.TODO(),
		"",                          // exchange
		constants.RabbitChannelName, // routing key
		false,                       // mandatory
		false,                       // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        m.ToBytes(),
			MessageId:   m.Sender + m.Receiver + uuid.New().String(),
		},
	)

	if err != nil {
		return err
	}

	return nil
}

func Consume() (<-chan amqp.Delivery, error) {
	msgs, err := channel.Consume(
		constants.RabbitChannelName, // queue
		"test",                      // consumer
		true,                        // auto-ack
		false,                       // exclusive
		false,                       // no-local
		true,                        // no-wait
		nil,                         // args
	)

	if err != nil {
		return nil, err
	}

	return msgs, nil
}

func (m *Message) ToBytes() []byte {
	b, err := json.Marshal(m.Message)
	if err != nil {
		return nil
	}

	return b
}
