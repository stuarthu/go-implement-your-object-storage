package rabbitmq

import (
	"encoding/json"
	"github.com/streadway/amqp"
)

type Queue interface {
	Publish(replyTo string, body interface{})
	Consume() <-chan amqp.Delivery
}

type RabbitMQ struct {
	channel  *amqp.Channel
	Name     string
	exchange string
}

func New(s string) *RabbitMQ {
	conn, er := amqp.Dial(s)
	if er != nil {
		panic(er)
	}

	ch, er := conn.Channel()
	if er != nil {
		panic(er)
	}

	q, er := ch.QueueDeclare(
		"",    // name
		false, // durable
		true,  // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if er != nil {
		panic(er)
	}

	mq := new(RabbitMQ)
	mq.channel = ch
	mq.Name = q.Name
	return mq
}

func (q *RabbitMQ) Bind(e string) {
	er := q.channel.ExchangeDeclare(
		e,
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)
	if er != nil {
		panic(er)
	}
	er = q.channel.QueueBind(
		q.Name, // queue name
		"",     // routing key
		e,      // exchange
		false,
		nil)
	if er != nil {
		panic(er)
	}
	q.exchange = e
}

func (q *RabbitMQ) Send(queue string, body interface{}) {
	str, e := json.Marshal(body)
	if e != nil {
		panic(e)
	}
	e = q.channel.Publish("",
		queue,
		false,
		false,
		amqp.Publishing{
			ReplyTo: q.Name,
			Body:    []byte(str),
		})
	if e != nil {
		panic(e)
	}
}

func (q *RabbitMQ) Publish(exchange string, body interface{}) {
	str, e := json.Marshal(body)
	if e != nil {
		panic(e)
	}
	e = q.channel.Publish(exchange,
		"",
		false,
		false,
		amqp.Publishing{
			ReplyTo: q.Name,
			Body:    []byte(str),
		})
	if e != nil {
		panic(e)
	}
}

func (q *RabbitMQ) Consume() <-chan amqp.Delivery {
	c, e := q.channel.Consume(q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if e != nil {
		panic(e)
	}
	return c
}

func (q *RabbitMQ) Close() {
	q.channel.Close()
}
