package rabbitmq

import (
	"fmt"
	"testing"

	"github.com/streadway/amqp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type TestChannel struct {
	*amqp.Channel
	mock.Mock
}

func (t *TestChannel) Publish(exchange string, key string, mandatory bool, immediate bool, msg amqp.Publishing) error {
	args := t.Called(exchange, key, mandatory, immediate, msg)
	return args.Error(0)
}

func (t *TestChannel) ExchangeDeclare(exchange, exchangeType string, durable, autoDelete, internal, noWait bool, args amqp.Table) error {
	arguments := t.Called(exchange, exchangeType, durable, autoDelete, internal, noWait, args)
	return arguments.Error(0)
}

func (t *TestChannel) QueueDeclare(queueName string, durable, autoDelete, exclusive, noWait bool, args amqp.Table) (amqp.Queue, error) {
	arguments := t.Called(queueName, durable, autoDelete, exclusive, noWait, args)
	return arguments.Get(0).(amqp.Queue), arguments.Error(1)
}

func (t *TestChannel) QueueBind(queueName, routingKey, exchangeName string, noWait bool, args amqp.Table) error {
	arguments := t.Called(queueName, routingKey, exchangeName, noWait, args)
	return arguments.Error(0)
}

type TestConnection struct {
	mock.Mock
	*amqp.Connection
}

func TestPush_SuccessfulPublish(t *testing.T) {
	// arrange
	amqpQueue := amqp.Queue{Name: "test.queue"}
	exchangeName := "test"
	routingKey := "work"

	testChannel := new(TestChannel)
	rabbitmq := rabbitMQ{
		channel:    testChannel,
		connection: &TestConnection{},
		ctag:       "test",
	}

	testChannel.On("ExchangeDeclare", exchangeName, "direct", true, false, false, false, mock.Anything).Return(nil)
	testChannel.On("QueueDeclare", amqpQueue.Name, true, false, false, false, mock.Anything).Return(amqpQueue, nil)
	testChannel.On("QueueBind", amqpQueue.Name, routingKey, exchangeName, false, mock.Anything).Return(nil)
	testChannel.On("Publish", exchangeName, routingKey, false, false, mock.Anything).Return(nil)

	queue, _ := rabbitmq.NewRabbitQueue(amqpQueue.Name, exchangeName, routingKey, nil)

	message := struct {
		Test  string
		Stuff string
	}{"this", "test"}

	// act
	err := queue.Publish(message)

	// assert
	assert.Nil(t, err)
}

func TestPush_FailedPublish(t *testing.T) {
	// arrange
	amqpQueue := amqp.Queue{Name: "test.queue"}
	exchangeName := "test"
	routingKey := "work"

	testChannel := new(TestChannel)
	rabbitmq := rabbitMQ{
		channel:    testChannel,
		connection: &TestConnection{},
		ctag:       "test",
	}

	testChannel.On("ExchangeDeclare", exchangeName, "direct", true, false, false, false, mock.Anything).Return(nil)
	testChannel.On("QueueDeclare", amqpQueue.Name, true, false, false, false, mock.Anything).Return(amqpQueue, nil)
	testChannel.On("QueueBind", amqpQueue.Name, routingKey, exchangeName, false, mock.Anything).Return(nil)
	testChannel.On("Publish", exchangeName, routingKey, false, false, mock.Anything).Return(fmt.Errorf("Failure"))

	queue, _ := rabbitmq.NewRabbitQueue(amqpQueue.Name, exchangeName, routingKey, nil)

	message := struct {
		Test  string
		Stuff string
	}{"this", "test"}

	// act
	err := queue.Publish(message)

	// assert
	assert.NotNil(t, err)
}
