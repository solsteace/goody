package amqp

import (
	"fmt"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

// Inspiration: https://github.com/streadway/amqp/issues/403
type publisher struct {
	reconnectDelay int
	conn           *amqp091.Connection
	url            string
	err            chan error
	amqpErr        chan *amqp091.Error

	channels map[string]*channel
}

func NewPublisher(url string, reconnectDelay int) publisher {
	return publisher{
		reconnectDelay: reconnectDelay,
		url:            url,
		err:            make(chan error),
		amqpErr:        make(chan *amqp091.Error),
		channels:       make(map[string]*channel),
	}
}

func (p publisher) Initiate() {
	go (func() { // Reconnection
		for {
			select {
			case <-p.err:
				p.reconnect()
			case <-p.amqpErr:
				p.reconnect()
			}
		}
	})()
	p.connect() // Connection
}

func (p *publisher) Track(c *channel, name string) {
	if _, ok := p.channels[name]; ok {
		return
	}
	p.channels[name] = c
}

// TODO: smarter reconnection rather than total reconnecttion whenever a failure
// during initating a component happened
func (p *publisher) connect() {
	conn, err := amqp091.Dial(p.url)
	if err != nil {
		p.err <- err
		return
	}

	conn.NotifyClose(p.amqpErr)
	p.conn = conn

	for name, c := range p.channels {
		if err := c.initiate(conn); err != nil {
			p.err <- err
		}
		fmt.Printf("Channel OK: %s\n", name)
	}
}

func (p publisher) reconnect() {
	time.Sleep(time.Second * time.Duration(p.reconnectDelay))
	p.connect()
}
