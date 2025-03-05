package main

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/thebigyovadiaz/rabbitmq-hello-world/src/util"
)

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672")
	util.LogFailOnError(err, "Failed to connect to RabbitMQ")
	util.LogSuccessful("Connected to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	util.LogFailOnError(err, "Failed to open a channel")
	util.LogSuccessful("Channel open successfully")
	defer ch.Close()

	qD, err := ch.QueueDeclare(
		"hello",
		false,
		false,
		false,
		false,
		nil,
	)
	util.LogFailOnError(err, "Failed to declare a queue")

	messages, err := ch.Consume(
		qD.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	util.LogFailOnError(err, "Failed to register a consumer")

	var forever chan struct{}

	go func() {
		for d := range messages {
			util.LogSuccessful(fmt.Sprintf("Received a message: %s", d.Body))
		}
	}()

	util.LogSuccessful(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
