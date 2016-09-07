package rabbitmq

import (
	"github.com/streadway/amqp"
)

// AmqpChannel - Interface definition of streadway amqp channel
// this is used to allow for swapping out implementations during
// testing
type AmqpChannel interface {
	Ack(uint64, bool) error
	Cancel(string, bool) error
	Close() error
	Confirm(bool) error
	Consume(string, string, bool, bool, bool, bool, amqp.Table) (<-chan amqp.Delivery, error)
	ExchangeBind(string, string, string, bool, amqp.Table) error
	ExchangeDeclare(string, string, bool, bool, bool, bool, amqp.Table) error
	ExchangeDeclarePassive(string, string, bool, bool, bool, bool, amqp.Table) error
	ExchangeDelete(string, bool, bool) error
	ExchangeUnbind(string, string, string, bool, amqp.Table) error
	Flow(bool) error
	Get(string, bool) (amqp.Delivery, bool, error)
	Nack(uint64, bool, bool) error
	NotifyCancel(chan string) chan string
	NotifyClose(chan *amqp.Error) chan *amqp.Error
	NotifyConfirm(chan uint64, chan uint64) (chan uint64, chan uint64)
	NotifyFlow(chan bool) chan bool
	NotifyPublish(chan amqp.Confirmation) chan amqp.Confirmation
	NotifyReturn(chan amqp.Return) chan amqp.Return
	Publish(string, string, bool, bool, amqp.Publishing) error
	Qos(int, int, bool) error
	QueueBind(string, string, string, bool, amqp.Table) error
	QueueDeclare(string, bool, bool, bool, bool, amqp.Table) (amqp.Queue, error)
	QueueDeclarePassive(string, bool, bool, bool, bool, amqp.Table) (amqp.Queue, error)
	QueueDelete(string, bool, bool, bool) (int, error)
	QueueInspect(string) (amqp.Queue, error)
	QueuePurge(string, bool) (int, error)
	QueueUnbind(string, string, string, amqp.Table) error
	Recover(bool) error
	Reject(uint64, bool) error
	Tx() error
	TxCommit() error
	TxRollback() error
}
