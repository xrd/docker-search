package main

import (
	"testing"
	"path"
	"os"
)

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

func TestFilters( t* testing.T ) {
	c := new(Client)
	c.LoadConfig( path.Join( ".", "test", "fixtures", DEFAULT_CONFIG_FILE ) )
	c.Query( "ffmpeg" )
	results := c.Filter( []string{"libavcodec=src", "python"} )
	if "" != results["link/ffmpeg-built"] {
		t.Errorf( "Unable to find ffmpeg library with python and libavcodec" )
	}
	
}

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

func TestFilter( t* testing.T ) {
	file, _ := os.Open( path.Join( "test", "fixtures", "bfirst_ffmpeg_dockerfile.txt" ) )
	if nil != file {
		
	} else {
		t.Errorf( "Nope, unable to open!" )
	}
}
