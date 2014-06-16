package main

import (
	"net/http"
	"io/ioutil"
	"fmt"
	"github.com/BurntSushi/toml"
)

type Client struct {
	config *Config
	dockerfiles map[string]string
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

func (c* Client) Search() string { 
	return "Hi there"
}

func (c* Client) Query( term string ) string {

	// GET /v1/search?q=search_term HTTP/1.1
	// Host: example.com
	// Accept: application/json
	client := &http.Client{}
	url := c.config.Host + c.config.Endpoint + "?term=ffmpeg"
	req, _ := http.NewRequest( "GET", url, nil )
	req.Header.Add( "Accept", "application/json")
	req.Header.Add( "User-Agent", "Docker-Client/1.0.0" )
	req.Header.Add( "X-Registry-Auth", "eyJhdXRoIjoiIiwiZW1haWwiOiIifQ==" )
	if resp, err := client.Do(req); nil != err {
		fmt.Println( "Error: ", err )
		return ""
	} else {
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		return string(body)
	}
}
