package main

import (
	"context"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/thebigyovadiaz/rabbitmq-hello-world/src/util"
	"log"
	"time"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func connSuccessful(msg string) {
	log.Println(msg)
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	util.LogFailOnError(err, "Failed to connect to RabbitMQ")
	util.LogSuccessful("RabbitMQ connect successfully")
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
	failOnError(err, "Failed to declare a queue")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	testCases := []string{
		"Hello World with RabbitMQ",
		"Hello World with RabbitMQ",
		"Hello World with RabbitMQ",
		"Hello World with RabbitMQ",
		"Hello World with RabbitMQ",
	}

	for k, v := range testCases {
		body := fmt.Sprintf("%s %d", v, k+1)
		err = ch.PublishWithContext(
			ctx,
			"",
			qD.Name,
			false,
			false,
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(body),
			},
		)

		failOnError(err, "Failed to publish a message")
		connSuccessful(fmt.Sprintf(" [x] Sent %s", body))
	}

	connSuccessful("Messages sent successfully")
}
