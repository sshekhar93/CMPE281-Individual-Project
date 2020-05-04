package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	entity "github.com/cmpe281-sshekhar93/bitly/control-panel/linkEntity"
	service "github.com/cmpe281-sshekhar93/bitly/control-panel/service"
	"github.com/gorilla/mux"
)

//LinkController : controller interface
type LinkController interface {
	CreateLink(resp http.ResponseWriter, req *http.Request)
	GetLink(resp http.ResponseWriter, req *http.Request)
	GetLinks(resp http.ResponseWriter, req *http.Request)
	UpdateLink(resp http.ResponseWriter, req *http.Request)
	DeleteLink(resp http.ResponseWriter, req *http.Request)
	//RedirectLink(resp http.ResponseWriter, req *http.Request)
}

// controller : Implements route functions
type controller struct{}

//NewController : Constructor to get contoller instance
func NewController() LinkController {
	return &controller{}
}

var linkService service.Service = service.NewService()

//CreateLink : Serves /create POST methodroute and creates a new short link for URL
func (*controller) CreateLink(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "application/json")
	var linkData entity.LinkData
	err := json.NewDecoder(req.Body).Decode(&linkData)
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(`{"error": "Error decoding the request failed"}`))
		return
	}
	fmt.Println("Calling create of service")
	err = linkService.Create(&linkData)
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(`{"error": "Error: failed to generate short url link"}`))
		return
	}
	result, err := json.Marshal(linkData)
	resp.WriteHeader(http.StatusOK)
	resp.Write(result)
}

//GetLink : Serves /link/{id} route GET method and returns all the details for the provided id
func (*controller) GetLink(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "application/json")
	var linkData entity.LinkData
	vars := mux.Vars(req)
	id := vars["id"]

	err := linkService.GetLink(&linkData, id)
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(`{"error": "Error: failed to get shortlink details"}`))
		return
	}
	result, err := json.Marshal(linkData)
	resp.WriteHeader(http.StatusOK)
	resp.Write(result)
}

//GetLink : Serves /links route GET method and returns all the stred short Links
func (*controller) GetLinks(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "application/json")
	linkDatas, err := linkService.GetLinks()
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(`{"error": "Error: failed to get shortlink details"}`))
		return
	}
	result, err := json.Marshal(linkDatas)
	resp.WriteHeader(http.StatusOK)
	resp.Write(result)
}

//UpdateLink : Serves /link/{id} route PUT method and updates the uri for given id
func (*controller) UpdateLink(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "application/json")
	var linkData entity.LinkData
	err := json.NewDecoder(req.Body).Decode(&linkData)
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(`{"error": "Error decoding the request failed"}`))
		return
	}
	err = linkService.Update(&linkData)
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(`{"error": "Error: failed to generate short url link"}`))
		return
	}
	result, err := json.Marshal(linkData)
	resp.WriteHeader(http.StatusOK)
	resp.Write(result)
}

//DeleteLink : Serves /link/{id} route DELETE method and delete the data for given id
func (*controller) DeleteLink(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(req)
	id := vars["id"]

	err := linkService.Delete(id)
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(`{"error": "Error: failed to get shortlink details"}`))
		return
	}
	resp.WriteHeader(http.StatusOK)
	resp.Write([]byte("ID: " + id + " deleted successfully"))
}

//RedirectLink : Serves /{shortLink} route GET method and redirects to actual URL
/*func (*controller) RedirectLink(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "application/json")
	fmt.Println("RedirectLink(): ")
	var linkData entity.LinkData
	shortLink := "http://" + req.Host + req.URL.Path
	fmt.Println("RedirectLink(): ", linkData, shortLink)
	err := linkService.GetLinkRedirect(&linkData, shortLink)
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(`{"error": "Error: failed to get shortlink details"}`))
		return
	}
	fmt.Println(linkData)
	http.Redirect(resp, req, linkData.Uri, http.StatusSeeOther)
}*/
