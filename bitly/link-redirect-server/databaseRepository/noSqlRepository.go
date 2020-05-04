package repository

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	entity "github.com/cmpe281-sshekhar93/bitly/link-redirect-server/linkEntity"
)

type noSqlRepo struct{}

func NewNoSqlRepository() Repository {
	return &noSqlRepo{}
}

func (*noSqlRepo) FETCHREDIRECT(linkData *entity.LinkData, shortLink string) error {
	slice := strings.Split(shortLink, "/")
	key := slice[len(slice)-1]
	fmt.Println("key: " + key)
	resp, err := http.Get("http://bitly-nosql-net-elb-29243d5b0b215828.elb.us-west-2.amazonaws.com:9090/api/" + key) //TODO : add actual load balancer address
	if err != nil {
		fmt.Println("linkCreateRoutine(): nosql get failed" + err.Error())
		if resp != nil {
			fmt.Println("linkCreateRoutine(): nosql get failed", resp.Body)
			resp.Body.Close()
		}
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		if resp != nil {
			resp.Body.Close()
		}
	}
	err = json.Unmarshal(body, &linkData)
	return err
}

func (*noSqlRepo) INSERT(linkData *entity.LinkData) error {
	return nil
}
func (*noSqlRepo) FETCH(linkData *entity.LinkData, id string) error {
	return nil
}
func (*noSqlRepo) FETCHALL() ([]entity.LinkData, error) {
	return []entity.LinkData{}, nil
}
func (*noSqlRepo) UPDATE(linkData *entity.LinkData) error {
	return nil
}
func (*noSqlRepo) DELETE(id string) error {
	return nil
}
