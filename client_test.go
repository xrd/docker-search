package main

import "testing"

func TestSomething( t *testing.T ) {
	c := new(Client)
	files := make( map[string]string )
	files["foobar"] = "Hi there"
	c.Load( files )
	// c.Query( "hi" )
	if result := c.Filter( []string{"Hi"} ); "" == result {
		t.Errorf( "Error!" )
	}

}
