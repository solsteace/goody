package amqp

import (
	"github.com/rabbitmq/amqp091-go"
)

type channel struct {
	Conn   *amqp091.Channel
	queues map[string]*queue
}

func NewChannel() channel {
	channel := channel{queues: make(map[string]*queue)}
	return channel
}

func (ce *channel) Track(q *queue) {
	if _, ok := ce.queues[q.cfgName]; ok {
		return
	}
	ce.queues[q.cfgName] = q
}

func (ce *channel) initiate(conn *amqp091.Connection) error {
	c, err := conn.Channel()
	if err != nil {
		return err
	}

	ce.Conn = c
	for _, queue := range ce.queues {
		if err := queue.initiate(ce.Conn); err != nil {
			return err
		}
	}
	return nil
}
