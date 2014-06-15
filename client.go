package main

import (
	"net/http"
	"io/ioutil"
)

type Docker struct {
	terms []string
}

var dockerfiles map[string]string

func Load( files map[string]string ) {
	dockerfiles = files
}

func Search( term string ) string {

	// GET /v1/search?q=search_term HTTP/1.1
	// Host: example.com
	// Accept: application/json
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "http://hub.docker.io/v1/search?q=ffmpeg", nil)
	// ...
	req.Header.Add("Accept", "application/json")
	resp, _ := client.Do(req)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body)
}
