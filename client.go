package main

import (
	"net/http"
	"io/ioutil"
	// "fmt"
	// "os"
	"github.com/BurntSushi/toml"
	"strings"
	"encoding/json"
	"log"
	"html"
)

type Client struct {
	config *Config
	dockerfiles map[string]string
	Images []DockerImage
	Results []DockerImage
	Verbose bool
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

func (c* Client) log( msg ...string ) {
	logit( c.Verbose, msg... )
}

func (c* Client) LoadConfig( path string ) bool {
	rv := true
	var conf Config
	if _, err := toml.DecodeFile( path, &conf ); err != nil {
		// handle error
		log.Fatal( err )
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

func (c* Client) grabDockerfile( ci chan<- Tuple, name string ) {
	// Raw link: https://registry.hub.docker.com/u/bfirsh/ffmpeg/dockerfile/raw
	client := &http.Client{}
	url := "https://registry.hub.docker.com/u/" + name + "/dockerfile/raw"
	req, _ := http.NewRequest( "GET", url, nil )
	if resp, err := client.Do(req); nil != err {
		log.Fatal( err )
	} else {
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		ci <- Tuple{name,string(body)}
	}
}

func (c* Client) processDockerfile( ci <-chan Tuple ) {
	tuple := <- ci
	// Apply it to the correct result
	for i,image := range c.Images {
		if tuple.Name == image.Name {
			c.log(  "Got dockerfile for: " + tuple.Name ) 
			c.Images[i].Dockerfile = strings.TrimSpace( html.UnescapeString( tuple.Dockerfile ) )
		}
	}
	
}

func (c* Client) Annotate() {

	// Grab a bunch of Dockerfiles, and then process them
	count := 0
	ci := make(chan Tuple, 4 )
	for _, image := range c.Images {
		c.log( "Annotating dockerfile for " + image.Name )
		go c.grabDockerfile( ci, image.Name )
		count++
		
	}
	for count > 0 {
		c.processDockerfile( ci )
		count--
	}
	c.log( "Finished annotation of dockerfiles" )
}

func (c* Client) Filter( filters []string ) {
	
	results := make( map[string]DockerImage )

	out := ""
	for _, e := range filters {
		out += " " + e
	}
	c.log( "Filters: " + out )

	if 0 < len( filters ) {
		c.log( "Filtering dockerfiles" )
		
		for _, filter := range filters {
			
			td := ProcessFilter( filter )
			for _, image := range c.Images {
				if -1 != strings.Index( image.Dockerfile, td.Target ) {
					c.log( "Found match inside Dockerfile" )
					results[image.Name] = image
				}
			}
		}
		
		// Set them to the results
		c.Results = []DockerImage{}
		for _,v := range results {
			c.log( "Adding result to results: " + v.Name )
			c.Results = append( c.Results, v )
		}
	} else {
		c.Results = c.Images
	}
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

type DockerResults struct {
	Results []DockerImage `json:"results"`
	Query string `json:"query"`
}

type DockerImage struct {
        Description string
        IsOfficial bool `json:"is_official"`
        IsTrusted bool `json:"is_trusted"`
        Name string
        StarCount int `json:"star_count"`
	Dockerfile string `json:"dockerfile"`
}

func (c* Client) Query( term string ) bool {

	rv := true
	// GET /v1/search?q=search_term HTTP/1.1
	// Host: example.com
	// Accept: application/json
	client := &http.Client{}
	url := c.config.Host + c.config.Endpoint + "?q=" + term
	req, _ := http.NewRequest( "GET", url, nil )
	req.Header.Add( "Accept", "application/json")
	req.Header.Add( "User-Agent", "Docker-Client/1.0.0" )
	if resp, err := client.Do(req); nil != err {
		log.Fatal( err )
		rv = false

	} else {
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		var res DockerResults
		err := json.Unmarshal(body, &res)
		if nil == err {
			c.Images = res.Results
			c.Results = res.Results
		} else {
			log.Fatal( err )
		}
	}
	return rv
}
