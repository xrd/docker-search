package main

import (
	"net/http"
	"io/ioutil"
	"fmt"
	"github.com/BurntSushi/toml"
//	"strings"
//	"encoding/json"
)

type Client struct {
	config *Config
	dockerfiles map[string]string
	results string
}

type Docker struct {
	terms []string
}

func (c *Client) Load( files map[string]string ) {
	c.dockerfiles = files
}

type Config struct {
	Host string
	Endpoint string
	UpdateCheck bool
}

func (c* Client) LoadConfig( path string ) bool {
	rv := true
	var conf Config
	if _, err := toml.DecodeFile( path, &conf ); err != nil {
		// handle error
		fmt.Println( "Error: ", err )
		rv = false
	} else {
		c.config = &conf
	}
	return rv

}


type Tuple struct {
	name string
	dockerfile string
	
}

func (c* Client) grabDockerfile( ci chan Tuple, name string ) {
	// Raw link: https://registry.hub.docker.com/u/bfirsh/ffmpeg/dockerfile/raw
	client := &http.Client{}
	url := "https://registry.hub.docker.com/u/bfirsh/" + name + "/raw"
	req, _ := http.NewRequest( "GET", url, nil )
	if resp, err := client.Do(req); nil != err {
		fmt.Println( "Error: ", err )
	} else {
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		ci <- Tuple{name,string(body)}
	}
}

func (c* Client) Filter( items []string ) map[string]string { 

	// Grab a bunch of Dockerfiles, and then process them
	
	//count := 0
	//	ci := make(chan Tuple)
	// for i, e := range c.results {
	// 	// go c.grabDockerfile( ci, e["name"] )
	// 	ci <- Tuple{e,"something"}
	// 	count++
		
	// }
	// found := make( map[string]string )
	// for count > 0 {
	// 	tuple <- ci
	// 	found[tuple.name] = tuple.dockerfile
	// }

	// Process it all
	results := make( map[string]string )
	return results
}

func (c* Client) Query( term string ) {

	// GET /v1/search?q=search_term HTTP/1.1
	// Host: example.com
	// Accept: application/json
	client := &http.Client{}
	url := c.config.Host + c.config.Endpoint + "?term=" + term
	req, _ := http.NewRequest( "GET", url, nil )
	req.Header.Add( "Accept", "application/json")
	req.Header.Add( "User-Agent", "Docker-Client/1.0.0" )
	// req.Header.Add( "X-Registry-Auth", "eyJhdXRoIjoiIiwiZW1haWwiOiIifQ==" )
	if resp, err := client.Do(req); nil != err {
		fmt.Println( "Error: ", err )

	} else {
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		//var f interface{}
		//err := json.Unmarshal(body, &f)
		c.results = string(body)
	}
	fmt.Println( "Results: " + c.results )
}
