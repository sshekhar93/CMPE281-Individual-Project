package service

import (
	entity "github.com/cmpe281-sshekhar93/bitly/trend-server/linkEntity"
	"github.com/cmpe281-sshekhar93/bitly/trend-server/messanger"
)

type Service interface {
	Consume(queueName string, c chan entity.LinkTrendData) error
}

var trendMessanger messanger.Messanger = messanger.NewRabbitMessanger()

type trendService struct{}

func NewTrendService() Service {
	return &trendService{}
}

func (*trendService) Consume(queueName string, c chan entity.LinkTrendData) error {
	trendMessanger.CONSUME(queueName, c)
	return nil
}
