package amqp

import (
	"fmt"

	"github.com/rabbitmq/amqp091-go"
)

type queue struct {
	Conn          amqp091.Queue
	cfgName       string
	cfgDurable    bool
	cfgAutoDelete bool
	cfgExclusive  bool
	cfgNoWait     bool
	cfgArgs       amqp091.Table
}

func NewQueue(
	name string,
	durable,
	autoDelete,
	exclusive,
	noWait bool,
	args amqp091.Table,
) queue {
	return queue{
		cfgName:       name,
		cfgDurable:    durable,
		cfgAutoDelete: autoDelete,
		cfgExclusive:  exclusive,
		cfgNoWait:     noWait,
		cfgArgs:       args}
}

func (q *queue) initiate(channel *amqp091.Channel) error {
	queue, err := channel.QueueDeclare(
		q.cfgName,       // name
		q.cfgDurable,    // durable (on shutdown, should the queue be persisted?)
		q.cfgAutoDelete, // delete when unused (when the last consumer leaves, should the queue be deleted?)
		q.cfgExclusive,  // exclusive (When the disconnected, should the queue be deleted?)
		q.cfgNoWait,     // no-wait
		q.cfgArgs,       // arguments
	)
	if err != nil {
		return err
	}
	fmt.Printf("Queue OK: %s\n", q.cfgName)

	q.Conn = queue
	return nil
}
