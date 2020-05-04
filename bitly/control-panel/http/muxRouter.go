package router

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type muxRouter struct{}

var muxDispatch = mux.NewRouter()

// NewMuxRouter : constructor forr Mux router
func NewMuxRouter() Router {
	return &muxRouter{}
}

func (*muxRouter) GET(uri string, f func(resp http.ResponseWriter, req *http.Request)) {
	muxDispatch.HandleFunc(uri, f).Methods("GET")
}

func (*muxRouter) POST(uri string, f func(resp http.ResponseWriter, req *http.Request)) {
	muxDispatch.HandleFunc(uri, f).Methods("POST")
}

func (*muxRouter) PUT(uri string, f func(resp http.ResponseWriter, req *http.Request)) {
	muxDispatch.HandleFunc(uri, f).Methods("PUT")
}

func (*muxRouter) DELETE(uri string, f func(resp http.ResponseWriter, req *http.Request)) {
	muxDispatch.HandleFunc(uri, f).Methods("DELETE")
}

func (*muxRouter) SERVE(port string) {
	fmt.Printf("Mux rrouter started")
	log.Fatalln(http.ListenAndServe(port, muxDispatch))
}
