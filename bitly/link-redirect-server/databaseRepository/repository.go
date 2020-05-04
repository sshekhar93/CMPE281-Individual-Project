package repository

import entity "github.com/cmpe281-sshekhar93/bitly/link-redirect-server/linkEntity"

//Repository inteface to use different database drivers
type Repository interface {
	FETCHREDIRECT(linkData *entity.LinkData, shortLink string) error
	INSERT(linkData *entity.LinkData) error
	FETCH(linkData *entity.LinkData, id string) error
	FETCHALL() ([]entity.LinkData, error)
	UPDATE(linkData *entity.LinkData) error
	DELETE(id string) error
}
