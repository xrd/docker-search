package main

import (
	"testing"
	"path"
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
