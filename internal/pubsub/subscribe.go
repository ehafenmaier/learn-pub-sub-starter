package pubsub

import (
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
)

func SubscribeJSON[T any](
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	simpleQueueType QueueType, // an enum to represent "durable" or "transient"
	handler func(T),
) error {
	// Declare and bind the queue
	ch, q, err := DeclareAndBind(
		conn,
		exchange,
		queueName,
		key,
		simpleQueueType)
	if err != nil {
		return err
	}

	// Get delivery channel from the queue
	messages, err := ch.Consume(q.Name, "", false, false, false, false, nil)
	if err != nil {
		return err
	}

	// Start a goroutine to handle messages
	go func() {
		for msg := range messages {
			var val T
			err := json.Unmarshal(msg.Body, &val)
			if err != nil {
				msg.Nack(false, false)
				continue
			}

			handler(val)
			msg.Ack(false)
		}
	}()

	return nil
}
