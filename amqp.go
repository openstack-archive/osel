package main

/*

amqp - This file includes all of the logic necessary to interact with the amqp
library.  This is extrapolated out so that a AmqpInterface interface can be
passed to functions.  Doing this allows testing by mock classes to be created
that can be passed to functions.

Since this is a wrapper around the amqp library, this does not need testing.

*/

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

// AmqpActioner is an interface for an AmqpActions class.  Having
// this as an interface allows us to pass in a dummy class for testing that
// just returns mocked data.
type AmqpActioner interface {
	Connect() (<-chan amqp.Delivery, error)
}

// AmqpActions is a class that handles all interactions directly with Amqp.
// See the comment on AmqpActioner for rationale.
type AmqpActions struct {
	Incoming       *<-chan amqp.Delivery
	Options        AmqpOptions
	AmqpConnection *amqp.Connection
	AmqpChannel    *amqp.Channel
	NotifyError    chan *amqp.Error
}

// AmqpOptions is a class to convey all of the configurable options for the
// AmqpActions class.
type AmqpOptions struct {
	RabbitURI string
}

// Connect initiates the initial connection to the AMQP.
func (s *AmqpActions) Connect() (<-chan amqp.Delivery, chan *amqp.Error, error) {
	var err error

	s.AmqpConnection, err = amqp.Dial(s.Options.RabbitURI)
	if err != nil {
		return nil, nil, fmt.Errorf("Failed to connect to RabbitMQ: %s", err)
	}
	s.NotifyError = s.AmqpConnection.NotifyClose(make(chan *amqp.Error)) //error channel

	s.AmqpChannel, err = s.AmqpConnection.Channel()
	if err != nil {
		return nil, nil, fmt.Errorf("Failed to open a channel: %s", err)
	}

	amqpQueue, err := s.AmqpChannel.QueueDeclare(
		"notifications.info", // name
		false,                // durable
		false,                // delete when usused
		false,                // exclusive
		false,                // no-wait
		nil,                  // arguments
	)
	if err != nil {
		return nil, nil, fmt.Errorf("Failed to declare a queue: %s", err)
	}

	amqpIncoming, err := s.AmqpChannel.Consume(
		amqpQueue.Name, // queue
		"osel",         // consumer
		true,           // auto-ack
		false,          // exclusive
		false,          // no-local
		false,          // no-wait
		nil,            // args
	)
	if err != nil {
		return nil, nil, fmt.Errorf("Failed to register a consumer: %s", err)
	}
	s.Incoming = &amqpIncoming
	return amqpIncoming, s.NotifyError, nil
}

// Close closes connections
func (s AmqpActions) Close() {
	log.Println("Closing AMQP connection")
	s.AmqpConnection.Close()
	s.AmqpChannel.Close()
}
