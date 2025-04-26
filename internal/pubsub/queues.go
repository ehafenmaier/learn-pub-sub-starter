package pubsub

import amqp "github.com/rabbitmq/amqp091-go"

type QueueType int

const (
	DurableQueue QueueType = iota
	TransientQueue
)

func DeclareAndBind(
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	simpleQueueType QueueType, // an enum to represent "durable" or "transient"
) (*amqp.Channel, amqp.Queue, error) {
	// Create a channel
	ch, err := conn.Channel()
	if err != nil {
		return nil, amqp.Queue{}, err
	}

	// Set queue properties based on the simpleQueueType
	durable := false
	autoDelete := false
	exclusive := false

	switch simpleQueueType {
	case DurableQueue:
		durable = true
	case TransientQueue:
		autoDelete = true
		exclusive = true
	default:
	}

	// Declare the queue
	q, err := ch.QueueDeclare(
		queueName,
		durable,
		autoDelete,
		exclusive,
		false,
		nil,
	)
	if err != nil {
		return nil, amqp.Queue{}, err
	}

	// Bind the queue to the exchange with the specified routing key
	err = ch.QueueBind(
		q.Name,
		key,
		exchange,
		false,
		nil,
	)
	if err != nil {
		return nil, amqp.Queue{}, err
	}

	// Return the channel and queue
	return ch, q, nil
}
