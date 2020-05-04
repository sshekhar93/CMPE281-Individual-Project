package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gorilla/mux"

	entity "github.com/cmpe281-sshekhar93/bitly/trend-server/linkEntity"
	"github.com/cmpe281-sshekhar93/bitly/trend-server/service"
)

type Controller interface {
	GetLinkTrend(resp http.ResponseWriter, req *http.Request)
	QueueConsumer(queueList []string) error
}

var (
	trendService service.Service = service.NewTrendService()
	c            chan entity.LinkTrendData
)

type trendController struct{}

func NewTrendController() Controller {
	return &trendController{}
}

func linkHitRoutine(queue string, c chan entity.LinkTrendData) {
	fmt.Println("Inside Link Hit routine")
	for linkTrendData := range c {
		fmt.Println(linkTrendData)
		slice := strings.Split(linkTrendData.ShotLink, "/")
		key := slice[len(slice)-1]
		fmt.Println("key: " + key)
		resp, err := http.Get("http://bitly-nosql-net-elb-29243d5b0b215828.elb.us-west-2.amazonaws.com:9090/api/" + key) //TODO : add actual load balancer address
		if err != nil {
			fmt.Println("linkCreateRoutine(): nosql get failed" + err.Error())
			//c <- linkTrendData
			if resp != nil {
				fmt.Println("linkCreateRoutine(): nosql get failed", resp.Body)
				resp.Body.Close()
			}
			continue
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			//c <- linkTrendData
			if resp != nil {
				resp.Body.Close()
			}
			continue
		}
		err = json.Unmarshal(body, &linkTrendData)
		if err != nil {
			//c <- linkTrendData
			if resp != nil {
				resp.Body.Close()
			}
			continue
		}
		linkTrendData.Count++
		jsonData, err := json.Marshal(linkTrendData)
		if err != nil {
			//c <- linkTrendData
			if resp != nil {
				resp.Body.Close()
			}
			continue
		}
		client := &http.Client{}
		fmt.Println("Sending PUT request with update count: ", linkTrendData.Count)
		request, err := http.NewRequest(http.MethodPut, "http://bitly-nosql-net-elb-29243d5b0b215828.elb.us-west-2.amazonaws.com:9090/api/"+key, bytes.NewBuffer(jsonData)) //TODO : add actual load balancer address
		if err != nil {
			fmt.Println("linkCreateRoutine(): nosql put failed" + err.Error())
			//c <- linkTrendData
			continue
		}
		request.Header.Set("Content-Type", "application/json; charset=utf-8")
		resp, err = client.Do(request)
		if (err != nil) || (resp.StatusCode != http.StatusOK) {
			//c <- linkTrendData
			if resp != nil {
				resp.Body.Close()
			}
			continue
		}
	}
}

func linkCreateRoutine(queue string, c chan entity.LinkTrendData) {
	for linkTrendData := range c {
		fmt.Println(linkTrendData)
		jsonData, err := json.Marshal(linkTrendData)
		slice := strings.Split(linkTrendData.ShotLink, "/")
		key := slice[len(slice)-1]
		fmt.Println("key: " + key)
		if err != nil {
			//c <- linkTrendData
			continue
		}
		resp, err := http.Post("http://bitly-nosql-net-elb-29243d5b0b215828.elb.us-west-2.amazonaws.com:9090/api/"+key, "application/json", bytes.NewBuffer(jsonData)) //TODO : add actual load balancer address
		if (err != nil) || (resp.StatusCode != http.StatusOK) {
			fmt.Println("linkCreateRoutine(): nosql post failed" + err.Error())
			//c <- linkTrendData
			if resp != nil {
				fmt.Println("linkCreateRoutine(): nosql post failed", resp.Body)
				resp.Body.Close()
			}
			continue
		}
	}
}

func (*trendController) QueueConsumer(queueList []string) error {
	for _, queue := range queueList {
		c = make(chan entity.LinkTrendData)
		trendService.Consume(queue, c)
		switch queue {
		case "linkHit":
			fmt.Println("Starting Link Hit ROUTINE")
			go linkHitRoutine(queue, c)
		case "linkCreate":
			fmt.Println("Starting Link create ROUTINE")
			go linkCreateRoutine(queue, c)
		}
	}
	return nil
}

func (*trendController) GetLinkTrend(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "application/json")
	linkTrendData := entity.LinkTrendData{}
	vars := mux.Vars(req)
	shortLink := vars["shortLink"]
	r, err := http.Get("Load Balancer Address" + shortLink) //TODO : configure actual load balancer address
	if err != nil {
		if resp != nil {
			defer r.Body.Close()
		}
		return
	}
	fmt.Println("before readiing response body")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		defer r.Body.Close()
		return
	}
	noSQLTrendData := make(map[string]entity.LinkTrendData)
	err = json.Unmarshal(body, &noSQLTrendData)
	linkTrendData = noSQLTrendData[shortLink]
	jsonData, err := json.Marshal(linkTrendData)
	resp.WriteHeader(http.StatusOK)
	resp.Write(jsonData)
}
