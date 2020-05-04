package main

import (
	"net/http"

	"github.com/cmpe281-sshekhar93/bitly/trend-server/controller"
	"github.com/cmpe281-sshekhar93/bitly/trend-server/messanger"
	"github.com/gorilla/mux"
)

var rabbitMessanger messanger.Messanger = messanger.NewRabbitMessanger()
var trendController controller.Controller = controller.NewTrendController()
var queueList []string
var port string = ":6000"

// var readHit chan<- entity.LinkTrendData
// var readCreate chan<- entity.LinkTrendData

func main() {
	queueList = append(queueList, "linkHit", "linkCreate")
	trendController.QueueConsumer(queueList)
	router := mux.NewRouter()
	router.HandleFunc("/", func(resp http.ResponseWriter, req *http.Request) {
		resp.Header().Set("Content-Type", "application/json")
		resp.WriteHeader(http.StatusOK)
		resp.Write([]byte("TS up and running"))
	}).Methods("GET")
	router.HandleFunc("/linkTrend/{shortLink}", trendController.GetLinkTrend).Methods("GET")
	http.ListenAndServe(port, router)
}
