package rabbitmq

import (
	"net"

	"github.com/streadway/amqp"
)

// AmqpConnection - Interface definition of streadway amqp connection
// this is used to allow for swapping out implementations during
// testing
type AmqpConnection interface {
	Channel() (*amqp.Channel, error)
	Close() error
	LocalAddr() net.Addr
	NotifyBlocked(c chan amqp.Blocking) chan amqp.Blocking
	NotifyClose(c chan *amqp.Error) chan *amqp.Error
}
