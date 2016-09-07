package rabbitmq

import (
	"encoding/json"

	"github.com/streadway/amqp"
)

// RabbitMQ wrapper for interacting with RabbitMQ
type RabbitMQ interface {
	Shutdown()
	NewRabbitQueue(queueName, exchangeName, routingKey string, args amqp.Table) (RabbitQueue, error)
}

// RabbitMQ wrapper for interacting with rabbit mq
type rabbitMQ struct {
	channel    AmqpChannel
	connection AmqpConnection
	ctag       string
}

// NewRabbitMQ construct for creating a new RabbitMQ connection and interacting
// with it
func NewRabbitMQ(url string, ctag string) (RabbitMQ, error) {
	connection, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	channel, err := connection.Channel()
	if err != nil {
		return nil, err
	}

	return &rabbitMQ{
		channel:    channel,
		connection: connection,
		ctag:       ctag,
	}, nil
}

// Shutdown provides a clean method for shutting down the connection and channel
// associated with RabbitMQ
func (r *rabbitMQ) Shutdown() {
	r.channel.Cancel(r.ctag, true)
	r.connection.Close()
}

// RabbitQueue wrapper for interacting with a RabbitMQ queue
type RabbitQueue interface {
	QueueName() string
	Publish(interface{}) error
	Inspect() (amqp.Queue, error)
}

type rabbitQueue struct {
	queueName string
	exchange  string
	key       string
	channel   AmqpChannel
}

// NewRabbitQueue construct used to define an exchange and queue and bind them together
func (r *rabbitMQ) NewRabbitQueue(queueName, exchangeName, routingKey string, args amqp.Table) (RabbitQueue, error) {
	if err := r.channel.ExchangeDeclare(exchangeName, "direct", true, false, false, false, nil); err != nil {
		return nil, err
	}

	if _, err := r.channel.QueueDeclare(queueName, true, false, false, false, args); err != nil {
		return nil, err
	}

	if err := r.channel.QueueBind(queueName, routingKey, exchangeName, false, nil); err != nil {
		return nil, err
	}

	return &rabbitQueue{
		queueName: queueName,
		exchange:  exchangeName,
		key:       routingKey,
		channel:   r.channel,
	}, nil
}

func (r *rabbitQueue) QueueName() string {
	return r.queueName
}

func (r *rabbitQueue) Publish(msg interface{}) error {
	jsonMsg, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	return r.channel.Publish(r.exchange, r.key, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        jsonMsg,
	})
}

func (r *rabbitQueue) Inspect() (amqp.Queue, error) {
	return r.channel.QueueInspect(r.queueName)
}
