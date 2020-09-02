package main

import (
	"github.com/streadway/amqp"
)

func publishData(usrToken string, data []byte) {
	conn, ch := setupAMQP("amqp://knot:knot@localhost:5672/")
	defer conn.Close()
	defer ch.Close()

	err := ch.ExchangeDeclare(
		"data.published",
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to declare an exchange")

	err = ch.Publish(
		"data.published",
		"",
		false,
		false,
		amqp.Publishing{
			Headers: amqp.Table{"Authorization": usrToken,},
			ContentType: "application/json",
			Body: []byte(data),
		})
	failOnError(err, "failed to publish data")
}

func setupAMQP (amqpURL string) (*amqp.Connection, *amqp.Channel){
	conn, err := amqp.Dial(amqpURL)
	failOnError(err, "Failed to connect to RabbitMQ")

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	return conn, ch
}
