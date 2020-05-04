package messanger

import (
	"encoding/json"
	"fmt"

	entity "github.com/cmpe281-sshekhar93/bitly/link-redirect-server/linkEntity"
	"github.com/streadway/amqp"
)

type rabbitMqMessanger struct{}

var (
	conn *amqp.Connection
	ch   *amqp.Channel
	q    amqp.Queue
)

func amqpConnection() error {
	conn, err := amqp.Dial("amqp://bitly:bitly@10.0.1.212:5672")
	if err != nil {
		return err
	}
	ch, err = conn.Channel()
	if err != nil {
		return err
	}
	q, err = ch.QueueDeclare(
		"linkHit",
		true,
		false,
		false,
		false,
		nil)
	if err != nil {
		return err
	}
	return err
}

// NewRabbitMessanger :  Constructor returning messanger of Rabbit type
func NewRabbitMessanger() Messanger {
	amqpConnection()
	return &rabbitMqMessanger{}
}

func (*rabbitMqMessanger) PUBLISH(linkData *entity.LinkData) error {
	fmt.Println(linkData)
	jsonMessage, err := json.Marshal(linkData)
	msg := amqp.Publishing{
		ContentType: "application/json",
		Body:        jsonMessage}
	ch.Publish(
		"",
		q.Name,
		false,
		false,
		msg)
	return err
}
func (*rabbitMqMessanger) CONSUME(linkData *entity.LinkData) error {
	return nil
}
