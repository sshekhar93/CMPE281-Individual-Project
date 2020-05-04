package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var Base uint64 = 1

type Range struct {
	Min uint64
	Max uint64
}

func getRange(resp http.ResponseWriter, req *http.Request) {
	countRange := Range{Base, (Base + 99999)}
	Base = countRange.Max + 1
	jsonRange, err := json.Marshal(countRange)
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte("failed to convert to json"))
		return
	}
	resp.Header().Set("Content-Type", "application/json")
	resp.WriteHeader(http.StatusOK)
	resp.Write(jsonRange)
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/range", getRange).Methods("GET")
	log.Fatalln(http.ListenAndServe(":3000", router))
}
