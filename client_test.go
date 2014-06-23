package main

import (
	"testing"
	"path"
	"os"
	"log"
	"strings"
)

var jsonSample = `
{ "query": "ffmpeg", "results": 
[{"description":"","is_official":false,"is_trusted":true,"name":"cellofellow/ffmpeg","star_count":1}
,{"description":"","is_official":false,"is_trusted":true,"name":"bfirsh/ffmpeg","star_count":0}
,{"description":"","is_official":false,"is_trusted":true,"name":"robd/aws-ffmpeg","star_count":0}
,{"description":"","is_official":false,"is_trusted":false,"name":"miovision/ffmpeg","star_count":0}
,{"description":"","is_official":false,"is_trusted":false,"name":"paulbrennan/ffmpeg","star_count":0}
,{"description":"","is_official":false,"is_trusted":false,"name":"cmark/ubuntu-ffmpeg","star_count":0}
]}
`

type MockWebClient struct {
}

func (mwc* MockWebClient) Get( url string ) []byte {
	var rv []byte
	if -1 != strings.Index( url, "index.docker.io" ) {
		// Return the JSON
		rv = []byte(jsonSample)
	} else if -1 != strings.Index( url, "registry.hub" ) {
		// Only if bfirsh/ffmpeg
		if -1 != strings.Index( url, "bfirsh/ffmpeg" ) {
			rv = []byte(FfmpegDockerfile)
		} else {
			rv = []byte("FROM ubuntu")
		}
	} 
	return rv
}

func getTestClient() *Client {
	c := new(Client)
	c.Http = new(MockWebClient)
	c.LoadConfig( path.Join( ".", "test", "fixtures", DEFAULT_CONFIG_FILE ) )
	// c.Verbose = true
	LoadFfmpegDockerfile()
	return c
}

func TestFullCircle( t *testing.T ) {
	c := getTestClient()

	c.Query( "ffmpeg" )
	c.Annotate()
	c.Filter( []string{ "quantal" } )

	if !c.resultFound( "bfirsh/ffmpeg" ) {
		t.Errorf( "Unable to filter for quantal" )
	}
	if c.resultFound( "miovision/ffmpeg" ) {
		t.Errorf( "This should not pass" )
	}
	if c.resultFound( "nevergonna/happen" ) {
		t.Errorf( "Hmm, we should not see this pass" )
	}
}

func (c* Client) resultFound( key string ) bool {
	found := false
	for _, e := range c.Results {
		if key == e.Name {
			found = true
		}
	}
	return found

}

// func TestFilters( t* testing.T ) {
// 	c := getTestClient()
// 	c.Query( "ffmpeg" )
// 	filters := []string{"libavcodec=src", "python"}
// 	c.Filter( filters )
	
// 	if !c.resultFound( "link/ffmpeg-built" ) {
// 		t.Errorf( "Unable to find ffmpeg library with python and libavcodec" )
// 	}
	
// }

func TestProcessFilter( t* testing.T ) {
	td := ProcessFilter( "libavcodec" )
	if td.Target != "libavcodec" && 
		!td.Src && 
		"" == td.Version {
		t.Errorf( "Did not process filter correctly" )
	}
	
	td = ProcessFilter( "libavcodec:src" )
	if td.Target != "libavcodec" && 
		td.Src {
		t.Errorf( "Did not process filter correctly" )
	}
	
	td = ProcessFilter( "libavcodec:2.2" )
	if td.Target != "libavcodec" && 
		!td.Src && 
		"2.2" == td.Version {
		t.Errorf( "Did not process filter correctly" )
	}

}


var FfmpegDockerfile = ""

func makeFakeImages() []DockerImage {
	rv := []DockerImage{}
	image := new(DockerImage)
	image.Name = "testing/test"
	image.Dockerfile = FfmpegDockerfile
	rv = append( rv, *image )
	return rv
}


func LoadFfmpegDockerfile() { 
	file, _ := os.Open( path.Join( "test", "fixtures", "bfirst_ffmpeg_dockerfile.txt" ) )
	if nil != file {
		contents := make([]byte, 1024 )
		_, err := file.Read( contents )
		if nil == err { 
			FfmpegDockerfile = string(contents)
		}
	} else {
		log.Fatal( "Cannot load test dockerfile" )
	}
}



func TestFilter( t* testing.T ) {

	c := getTestClient()
	c.Images = makeFakeImages()
	c.Filter( []string{ "quantal" } )
	if !c.resultFound( "testing/test" ) {
		t.Errorf( "Unable to filter for quantal" )
	}
	if c.resultFound( "nevergonna/happen" ) {
		t.Errorf( "Hmm, we should not see this pass" )
	}

}

// func (c* Client) Query() {
// 	var res DockerResults
// 	err := json.Unmarshal( jsonSample, res )
// 	if nil == err {
// 		c.Images = res.Results
// 	}
// }

// func (c* Client) Annotate() {
// 	for i, e := range c.Images {
// 		c.Images[i].Dockerfile = FfmpegDockerfile
// 	}
// }
