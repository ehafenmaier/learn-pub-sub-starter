package main

import (
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
	"time"
)

func publishGameLog(ch *amqp.Channel, username, message string) error {
	exchange := routing.ExchangePerilTopic
	key := routing.GameLogSlug + "." + username
	log := routing.GameLog{
		Username:    username,
		CurrentTime: time.Now(),
		Message:     message,
	}

	err := pubsub.PublishGob(ch, exchange, key, log)
	if err != nil {
		return err
	}

	return nil
}
