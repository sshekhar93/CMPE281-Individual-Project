package entity

//LinkData : data structure used to send http response and also to send value to mysql
type LinkTrendData struct {
	Id       int    `json: "id"`
	ShotLink string `json: "shortlink"`
	Uri      string `json: "uri"`
	Count    int    `json: "hits"`
}
