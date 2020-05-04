package messanger

import entity "github.com/cmpe281-sshekhar93/bitly/trend-server/linkEntity"

//Messanger : interface for message broker service
type Messanger interface {
	PUBLISH(linkData *entity.LinkTrendData) error
	CONSUME(queueName string, c chan entity.LinkTrendData) error
}
