package main

import (
	"testing"
	"path"
	"os"
)

var json_sample = `
[{"description":"","is_official":false,"is_trusted":true,"name":"cellofellow/ffmpeg","star_count":1}
,{"description":"","is_official":false,"is_trusted":true,"name":"bfirsh/ffmpeg","star_count":0}
,{"description":"","is_official":false,"is_trusted":true,"name":"robd/aws-ffmpeg","star_count":0}
,{"description":"FFMpeg built from source (git://source.ffmpeg.org/ffmpeg)","is_official":false,"is_trusted":false,"name":"lmars/ffmpeg","star_count":0}
,{"description":"this has python devenv, a few other build tools, and the open source libavcodec from ffmpeg built from source.","is_official":false,"is_trusted":false,"name":"link/ffmpeg-built","star_count":0}
,{"description":"","is_official":false,"is_trusted":false,"name":"miovision/ffmpeg","star_count":0}
,{"description":"","is_official":false,"is_trusted":false,"name":"paulbrennan/ffmpeg","star_count":0}
,{"description":"","is_official":false,"is_trusted":false,"name":"cmark/ubuntu-ffmpeg","star_count":0}
,{"description":"","is_official":false,"is_trusted":false,"name":"mikehearn/ubuntu-ffmpeg","star_count":0}
,{"description":"","is_official":false,"is_trusted":false,"name":"cmark/ubuntu-ffmpeg-ssh","star_count":0}
,{"description":"","is_official":false,"is_trusted":false,"name":"cmark/ubuntu-14.04-ffmpeg-nfs","star_count":0}
,{"description":"Docker based FFMPEG - keeping the dependancy soup in a neat Docker container.\n\nInstall ffmpeg from ppa:jon-severinsson/ffmpeg on ubuntu 12.04 container.","is_official":false,"is_trusted":false,"name":"asachs/docker-ffmpeg","star_count":0}
]
`

// func abcTestSomething( t *testing.T ) {
// 	c := new(Client)
// 	c.Query( "ffmpeg" )
// 	files := make( map[string]string )
// 	files["foobar"] = "Hi there"
// 	c.Load( files )
// 	// c.Query( "hi" )
// 	if result := c.Filter( []string{"Hi"} ); "" == result {
// 		t.Errorf( "Error!" )
// 	}

// }

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
// 	c := new(Client)
// 	c.LoadConfig( path.Join( ".", "test", "fixtures", DEFAULT_CONFIG_FILE ) )
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
	LoadFfmpegDockerfile()
	image.Dockerfile = FfmpegDockerfile
	rv = append( rv, *image )
	return rv
}


func LoadFfmpegDockerfile() { 
	file, _ := os.Open( path.Join( "test", "fixtures", "bfirst_ffmpeg_dockerfile.txt" ) )
	if nil != file {
		var contents []byte
		file.Read( contents )
		FfmpegDockerfile = string(contents)
	}
}

func TestFilter( t* testing.T ) {
	LoadFfmpegDockerfile()
	c := new(Client)
	c.LoadConfig( path.Join( ".", "test", "fixtures", DEFAULT_CONFIG_FILE ) )
	c.Images = makeFakeImages()
	c.Filter( []string{ "quantal" } )
	if !c.resultFound( "testing/test" ) {
		t.Errorf( "Unable to filter for quantal" )
	}
	if c.resultFound( "nevergonna/happen" ) {
		t.Errorf( "Hmm, we should not see this pass" )
	}

}
