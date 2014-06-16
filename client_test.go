package main

import "testing"

func TestSomething( t *testing.T ) {
	c := new(Client)
	files := make( map[string]string )
	files["foobar"] = "Hi there"
	c.Load( files )
	// c.Query( "hi" )
	if result := c.Search(); "" == result {
		t.Errorf( "Error!" )
	}

}
