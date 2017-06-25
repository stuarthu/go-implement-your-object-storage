package rabbitmq

import (
	"encoding/json"
	"testing"
)

func TestPublish(t *testing.T) {
	q := New("amqp://localhost:5672")
	defer q.Close()
	q.Bind("test")

	q2 := New("amqp://localhost:5672")
	defer q2.Close()
	q2.Bind("test")

	q3 := New("amqp://localhost:5672")
	defer q3.Close()
	q3.Bind("test2")

	q4 := New("amqp://localhost:5672")
	defer q4.Close()

	expect := "test"
	q4.Publish("test2", "any")
	q4.Publish("test", expect)

	c := q.Consume()
	msg := <-c
	var actual interface{}
	err := json.Unmarshal(msg.Body, &actual)
	if err != nil {
		t.Error(err)
	}
	if actual != expect {
		t.Errorf("expected %s, actual %s", expect, actual)
	}
	if msg.ReplyTo != q4.Name {
		t.Error(msg)
	}

	c2 := q2.Consume()
	msg = <-c2
	err = json.Unmarshal(msg.Body, &actual)
	if err != nil {
		t.Error(err)
	}
	if actual != expect {
		t.Errorf("expected %s, actual %s", expect, actual)
	}
	if msg.ReplyTo != q4.Name {
		t.Error(msg)
	}
	q2.Send(msg.ReplyTo, "test3")
	c3 := q4.Consume()
	msg = <-c3
	if string(msg.Body) != `"test3"` {
		t.Error(string(msg.Body))
	}
}

func TestSend(t *testing.T) {
	q := New("amqp://localhost:5672")
	defer q.Close()

	q2 := New("amqp://localhost:5672")
	defer q2.Close()

	expect := "test"
	expect2 := "test2"
	q2.Send(q.Name, expect)
	q2.Send(q2.Name, expect2)

	c := q.Consume()
	msg := <-c
	var actual interface{}
	err := json.Unmarshal(msg.Body, &actual)
	if err != nil {
		t.Error(err)
	}
	if actual != expect {
		t.Errorf("expected %s, actual %s", expect, actual)
	}

	c2 := q2.Consume()
	msg = <-c2
	err = json.Unmarshal(msg.Body, &actual)
	if err != nil {
		t.Error(err)
	}
	if actual != expect2 {
		t.Errorf("expected %s, actual %s", expect2, actual)
	}
}
