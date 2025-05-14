package pubsub

import (
	"bytes"
	"context"
	"encoding/gob"
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
)

func PublishJSON[T any](ch *amqp.Channel, exchange, key string, val T) error {
	// Marshal the value to JSON
	body, err := json.Marshal(val)
	if err != nil {
		return err
	}

	// Publish the message to the specified exchange with the given routing key
	err = ch.PublishWithContext(
		context.Background(),
		exchange, // exchange
		key,      // routing key
		false,    // mandatory
		false,    // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	if err != nil {
		return err
	}

	return nil
}

func PublishGob[T any](ch *amqp.Channel, exchange, key string, val T) error {
	// Encode the value to Gob
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(val)
	if err != nil {
		return err
	}

	// Publish the message to the specified exchange with the given routing key
	err = ch.PublishWithContext(
		context.Background(),
		exchange, // exchange
		key,      // routing key
		false,    // mandatory
		false,    // immediate
		amqp.Publishing{
			ContentType: "application/gob",
			Body:        buf.Bytes(),
		})
	if err != nil {
		return err
	}

	return nil
}
