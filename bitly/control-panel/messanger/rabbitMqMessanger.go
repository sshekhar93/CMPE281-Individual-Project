package messanger

import (
	"encoding/json"
	"fmt"

	entity "github.com/cmpe281-sshekhar93/bitly/control-panel/linkEntity"
	"github.com/streadway/amqp"
)

type rabbitMqMessanger struct{}

var (
	conn *amqp.Connection
	ch   *amqp.Channel
	q    amqp.Queue
)

func amqpConnection() error {
	fmt.Println("amqpConnection(): connection to rabbitMq server")
	conn, err := amqp.Dial("amqp://bitly:bitly@10.0.1.212:5672")
	if err != nil {
		fmt.Println("amqpConnection(): connection failed with error: " + err.Error())
		return err
	}
	ch, err = conn.Channel()
	if err != nil {
		return err
	}
	q, err = ch.QueueDeclare(
		"linkCreate",
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
