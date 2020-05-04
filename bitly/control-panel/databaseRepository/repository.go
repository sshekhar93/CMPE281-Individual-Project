package repository

import entity "github.com/cmpe281-sshekhar93/bitly/control-panel/linkEntity"

//Repository inteface to use different database drivers
type Repository interface {
	INSERT(linkData *entity.LinkData) error
	FETCH(linkData *entity.LinkData, id string) error
	FETCHALL() ([]entity.LinkData, error)
	//FETCHREDIRECT(linkData *entity.LinkData, shortLink string) error
	UPDATE(linkData *entity.LinkData) error
	DELETE(id string) error
}
