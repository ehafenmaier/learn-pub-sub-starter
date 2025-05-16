package pubsub

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
)

type AckType int

const (
	Ack AckType = iota
	NackRequeue
	NackDiscard
)

func SubscribeJSON[T any](
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	simpleQueueType QueueType, // an enum to represent "durable" or "transient"
	handler func(T) AckType,
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

	// Set prefetch count to 10
	err = ch.Qos(10, 0, false)
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
				fmt.Printf("Failed to unmarshal json: %s\n", err)
				continue
			}

			ackType := handler(val)
			switch ackType {
			case Ack:
				msg.Ack(false)
			case NackRequeue:
				msg.Nack(false, true)
			case NackDiscard:
				msg.Nack(false, false)
			}
		}
	}()

	return nil
}

func SubscribeGob[T any](
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	simpleQueueType QueueType,
	handler func(T) AckType,
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

	// Set prefetch count to 10
	err = ch.Qos(10, 0, false)
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
			buf := bytes.NewBuffer(msg.Body)
			dec := gob.NewDecoder(buf)
			err := dec.Decode(&val)
			if err != nil {
				fmt.Printf("Failed to decode gob: %s\n", err)
				continue
			}

			ackType := handler(val)
			switch ackType {
			case Ack:
				msg.Ack(false)
			case NackRequeue:
				msg.Nack(false, true)
			case NackDiscard:
				msg.Nack(false, false)
			}
		}
	}()

	return nil
}
