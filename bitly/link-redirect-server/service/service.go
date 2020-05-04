package service

import (
	"fmt"

	repository "github.com/cmpe281-sshekhar93/bitly/link-redirect-server/databaseRepository"
	entity "github.com/cmpe281-sshekhar93/bitly/link-redirect-server/linkEntity"
	"github.com/cmpe281-sshekhar93/bitly/link-redirect-server/messanger"
)

type Service interface {
	GetLinkRedirect(linkData *entity.LinkData, id string) error
}

type linkService struct{}

func NewService() Service {
	return &linkService{}
}

type Range struct {
	Min int `json: "min"`
	Max int `json: "max"`
}

var (
	repo      repository.Repository = repository.NewNoSqlRepository() //("bitly", "bitly", "10.0.1.97", "3306", "bitly")
	msgr      messanger.Messanger   = messanger.NewRabbitMessanger()
	linkCache                       = make(map[string]string)
)

func (*linkService) GetLinkRedirect(linkData *entity.LinkData, shortLink string) error {
	fmt.Println("GetLinkRedirect(): ", linkData, shortLink)
	uri, ok := linkCache[shortLink]
	if ok == false {
		err := repo.FETCHREDIRECT(linkData, shortLink)
		if err != nil {
			return err
		}
		if len(linkCache) >= 10 {
			for short, _ := range linkCache {
				delete(linkCache, short)
				break
			}
		}
		linkCache[shortLink] = linkData.Uri
	} else {
		linkData.ShotLink = shortLink
		linkData.Uri = uri
	}
	return msgr.PUBLISH(linkData)
}
