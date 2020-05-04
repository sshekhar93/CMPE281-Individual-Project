package messanger

import (
	"encoding/json"
	"fmt"

	entity "github.com/cmpe281-sshekhar93/bitly/trend-server/linkEntity"
	"github.com/streadway/amqp"
)

var (
	conn *amqp.Connection
	ch   *amqp.Channel
	// chCreate *amqp.Channel
)

type rabbitMqMessanger struct{}

func amqpConnection(rabbitUrl string) error {
	fmt.Println("Messanger amqpConnection(): " + rabbitUrl)
	conn, err := amqp.Dial(rabbitUrl)
	if err != nil {
		fmt.Println("Messanger amqpConnection(): " + rabbitUrl + err.Error())
		return err
	}
	ch, err = conn.Channel()
	// chCreate, err = conn.Channel()
	//fmt.Println("Messanger amqpConnection(): " + rabbitUrl + err.Error())
	return err
}

func NewRabbitMessanger() Messanger {
	amqpConnection("amqp://bitly:bitly@10.0.1.212:5672")
	return &rabbitMqMessanger{}
}

func consumerRoutine(queueName string, c chan entity.LinkTrendData) {
	fmt.Println("Messanger consumerRoutine(): " + queueName)
	linkTrendData := entity.LinkTrendData{}
	// var msgs <-chan amqp.Delivery
	// if queueName == "linkHit" {
	queue, err := ch.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil)
	if err != nil {
		return
	}

	msgs, err := ch.Consume(
		queue.Name,
		"",
		true,
		false,
		false,
		false,
		nil)
	if err != nil {
		return
	}
	// }
	// if queueName == "linkCreate" {
	// 	queue, err := chCreate.QueueDeclare(
	// 		queueName,
	// 		true,
	// 		false,
	// 		false,
	// 		false,
	// 		nil)
	// 	if err != nil {
	// 		return
	// 	}

	// 	msgs, err = chCreate.Consume(
	// 		queue.Name,
	// 		"",
	// 		true,
	// 		false,
	// 		false,
	// 		false,
	// 		nil)
	// 	if err != nil {
	// 		return
	// 	}
	// }
	for msg := range msgs {
		json.Unmarshal(msg.Body, &linkTrendData)
		c <- linkTrendData
	}
}

func (*rabbitMqMessanger) PUBLISH(linkData *entity.LinkTrendData) error {
	return nil
}

func (*rabbitMqMessanger) CONSUME(queueName string, c chan entity.LinkTrendData) error {
	fmt.Println("Messanger CONSUME()" + queueName)
	go consumerRoutine(queueName, c)
	return nil
}
