package messanger

import (
	entity "github.com/cmpe281-sshekhar93/bitly/control-panel/linkEntity"
)

//Messanger : Interface to use any message broker service
type Messanger interface {
	PUBLISH(linkData *entity.LinkData) error
	CONSUME(linkData *entity.LinkData) error
}
