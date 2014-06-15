package main

import (
	"fmt"
	"flag"
)

func help() {
		doc := `

docker-search: A better way to search the docker registry

docker-search can be used to search for specific information inside 
the docker registry. Want a docker base image with a specific version of 
ruby? Want to know what build flags were used when compiling that ffmpeg
binary? docker-search can help with all of that. docker-search searches 
not only for the name but also peers into the Dockerfile used to build 
the base image, and even traces back through the image's ancestry. 

Examples:

docker-search --ruby=1.9  # Search for base images with ruby 1.9
docker-search --python=2.2,src # images with python 2.2 compiled from source 
docker-search --last=2 --dockerfile # print out dockerfile from 2nd to last cached search
docker-search --list # print out cached search list

Flags:

--dockerfile        # Print out full Dockerfile with results
--string            # Search for the string in the Dockerfile
--string=match      # Search for the string with a filter of match
--string=match,src  # Search for compilation via source for this package
--last=[=N]         # Use a cached query, optionally specify Nth item
--list              # Print cached query list

`
	fmt.Println( doc )

}

func main() {
	
	// var ip = flag.Int("flagname", 1234, "help message for flagname")
	flag.Parse()
	
	help()

	files := make( map[string]string )
	files["ubuntu"] = "FROM foobar"
	Load( files )
	result := Search( "foobar" )
	fmt.Println( "Search results: ", result )
}
