package main

import (
	"net/http"
	"io/ioutil"
	"fmt"
	"github.com/BurntSushi/toml"
	"strings"
	"encoding/json"
)

type Client struct {
	config *Config
	dockerfiles map[string]string
	images []DockerImage
	Results []DockerImage
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
	Name string
	Dockerfile string
	
}

func (c* Client) grabDockerfile( ci chan Tuple, name string ) {
	// Raw link: https://registry.hub.docker.com/u/bfirsh/ffmpeg/dockerfile/raw
	client := &http.Client{}
	url := "https://registry.hub.docker.com/u/" + name + "/dockerfile/raw"
	// fmt.Println( "Grabbing dockerfile for " + name + " with URL: " + url )
	req, _ := http.NewRequest( "GET", url, nil )
	if resp, err := client.Do(req); nil != err {
		fmt.Println( "Error: ", err )
	} else {
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		ci <- Tuple{name,string(body)}
	}
}

func (c* Client) Annotate() {

	// Grab a bunch of Dockerfiles, and then process them
	
	count := 0
	ci := make(chan Tuple)
	for _, image := range c.images {
		go c.grabDockerfile( ci, image.Name )
		count++
		
	}
	for count > 0 {
		tuple := <- ci
		// Apply it to the correct result
		for _,image := range c.images {
			if tuple.Name == image.Name {
				tuple.Dockerfile = tuple.Dockerfile
			}
		}
		count--
	}
}

func (c* Client) Filter( items []string ) {
	
}

type TargetDescription struct {
	Src bool
	Version string
	Target string
}

func ProcessFilter( needle string ) *TargetDescription {
	td := new(TargetDescription)
	td.Src = false
	td.Version = ""
	usingColon := strings.Index( needle, ":" ) 
	usingComma := strings.Index( needle, "," ) 
	if -1 !=usingColon || -1 != usingComma {
		// Split it up, using the correct delimiter
		delimiter := ":"
		if -1 != usingComma { 
			delimiter = ","
		} 
		pieces := strings.Split( needle, delimiter )
		needle = pieces[0]
		for _,e := range pieces[1:] {
			if "src" == e {
				td.Src = true
			} else {
				// assume it is the version
				td.Version = e
			}
		}
	}
	td.Target = needle
	return td
}

func Search( needle string, haystack string ) bool {
	// Do some post processing on the string
	td := ProcessFilter( needle )
	
	if -1 != strings.Index( haystack, td.Target ) {
		fmt.Println( "Found it!" )
	}
	return false
}

// [{"description":"","is_official":false,"is_trusted":true,"name":"cellofellow/ffmpeg","star_count":1}
// 	,{"description":"","is_official":false,"is_trusted":true,"name":"bfirsh/ffmpeg","star_count":0}

type DockerImage struct {
        Description string
        IsOfficial bool `json:"is_official"`
        IsTrusted bool `json:"is_trusted"`
        Name string
        StarCount int `json:"star_count"`
	Dockerfile string
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
	if resp, err := client.Do(req); nil != err {
		fmt.Println( "Error: ", err )

	} else {
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		// fmt.Println( "Body: " + string(body) )
		var images []DockerImage
		json.Unmarshal(body, &images)
		c.images = images
		c.Results = images
	}
}
