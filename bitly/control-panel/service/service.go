package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	repository "github.com/cmpe281-sshekhar93/bitly/control-panel/databaseRepository"
	hash "github.com/cmpe281-sshekhar93/bitly/control-panel/hashAlgoithm"
	entity "github.com/cmpe281-sshekhar93/bitly/control-panel/linkEntity"
	"github.com/cmpe281-sshekhar93/bitly/control-panel/messanger"
)

type Service interface {
	Create(linkData *entity.LinkData) error
	GetLink(linkData *entity.LinkData, id string) error
	GetLinks() ([]entity.LinkData, error)
	//GetLinkRedirect(linkData *entity.LinkData, id string) error
	Update(linkData *entity.LinkData) error
	Delete(id string) error
}

type linkService struct{}

func NewService() Service {
	fetchRange()
	return &linkService{}
}

type Range struct {
	Min int `json: "min"`
	Max int `json: "max"`
}

var (
	countRange Range                 = Range{}
	count      int                   = 0
	hasher     hash.HashALgorithm    = hash.NewHashBase62()
	repo       repository.Repository = repository.NewMysqlRepository("bitly", "bitly", "10.0.1.97", "3306", "bitly")
	msgr       messanger.Messanger   = messanger.NewRabbitMessanger()
)

const (
	kongIp        string = "http://54.213.80.31:8000"
	rangeServerIp string = "http://10.0.1.183:3000/range"
	redirectApi   string = "/lrs/"
)

func fetchRange() error {
	fmt.Println("Inside fetch() function")
	resp, err := http.Get(rangeServerIp)
	if err != nil {
		fmt.Println("Inside fetch() connection failed: " + err.Error())
		if resp != nil {
			defer resp.Body.Close()
		}
		return err
	}
	fmt.Println("before readiing response body")
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		defer resp.Body.Close()
		return err
	}
	// fmt.Println([]byte(body))
	err = json.Unmarshal(body, &countRange)
	if err != nil {
		defer resp.Body.Close()
		return err
	}
	fmt.Println(countRange)
	count = countRange.Min
	return err
}

func (*linkService) Create(linkData *entity.LinkData) error {
	var uriHash string
	if linkData.Uri == "" {
		return errors.New("No URI present in request")
	}
	// uriHash = hasher.Encode(10000)
	fmt.Println("count: %d rangeMax: %d", count, countRange.Max)
	if count > countRange.Max {
		err := fetchRange()
		if err != nil {
			return err
		}

	}
	uriHash = hasher.Encode(count)
	count++
	linkData.ShotLink = kongIp + redirectApi + uriHash
	fmt.Println("linkData.ShotLink: %s, uriHash: %s, count: %d", linkData.ShotLink, uriHash, count)
	err := repo.INSERT(linkData)
	if err != nil {
		return err
	}
	return msgr.PUBLISH(linkData)
}

// GetLink : fectches info for provided "id"
func (*linkService) GetLink(linkData *entity.LinkData, id string) error {
	err := repo.FETCH(linkData, id)
	return err
}

func (*linkService) GetLinks() ([]entity.LinkData, error) {
	linkDatas, err := repo.FETCHALL()
	return linkDatas, err
}

/*func (*linkService) GetLinkRedirect(linkData *entity.LinkData, shortLink string) error {
	fmt.Println("GetLinkRedirect(): ", linkData, shortLink)
	err := repo.FETCHREDIRECT(linkData, shortLink)
	return err
}*/

func (*linkService) Update(linkData *entity.LinkData) error {
	err := repo.UPDATE(linkData)
	return err
}

func (*linkService) Delete(id string) error {
	err := repo.DELETE(id)
	return err
}
