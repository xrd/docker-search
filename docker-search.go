package main

import (
	"fmt"
	"flag"
	"os"
	"os/user"
	"path"
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

docker-search -filter=ruby:1.9 ffmpeg # Search for base images with ffpmeg with ruby 1.9
docker-search -filter=quantal ffmpeg  # Search for base images base ubuntu quantal in the Dockerfile
docker-search -filter=libavcodec:2.2:src ffmpeg # images with python 2.2 compiled from source 
docker-search -dockerfile ffmpeg # print out full dockerfiles

Flags:

-generate-config   # Create a new configuration file from current defaults
-dockerfile[s]     # Print out full Dockerfile with each results
-filter=str        # Search for the string in the Dockerfile
-filter=str:src    # Search for the string with a filter of match
-info[rmation]     # Print as much information as possible (maintainer, etc.)
`
	fmt.Println( doc )

}

const DEFAULT_CONFIG_FILE = ".docker-search.toml"

func getHomeDir() string {
	usr, err := user.Current()
	if err != nil {
		fmt.Println( "Error: ", err )
	}
	return usr.HomeDir 
}

func getConfigFilePath() string {
	return path.Join( getHomeDir(), DEFAULT_CONFIG_FILE )
}

func generateDefaultConfiguration() {
	configPath := getConfigFilePath()
	
	// Generate configuration
	if file, err := os.Create( configPath ); nil == err {
		def := `
Host = "http://192.168.59.103:2375"
Endpoint = "/v1.12/images/search"
UpdateCheck = true
`
		file.Write( []byte(def) )
		fmt.Println( "Created new configuration file" )
	} else {
		fmt.Println( "Unable to create new configuration file at: ", configPath, err )
	}
}


type filters []string

type Value interface {
	String() string
	Set(string) error
}

func (fs *filters) String() string {
	values := ""
	for _, s := range *fs {
		values += s + ":"
	}
	return values
}
 
func (fs *filters) Set(value string) error {
	*fs = append(*fs, value)
	return nil
}

func main() {

	var filts filters
	var genCon = flag.Bool( "generate-config", false, "Generate a new default configuration file" )
	var printDockerfile = flag.Bool( "dockerfile", false, "Print out dockerfiles" )
	var printInfo = flag.Bool( "info", false, "Print out detailed information on the maintainer(s)" )
	flag.Var( &filts, "filter", "List of filters" )
	flag.Parse()

	if *genCon {
		generateDefaultConfiguration()
		help()
		// flag.PrintDefaults()
	} else {
		if "" == flag.Arg(0) {
			help()
		} else {
			c := new(Client)
			if c.LoadConfig( getConfigFilePath() ) {
				c.Query( flag.Arg(0) )
				c.Annotate()
				if flag.NArg() > 1 {
					c.Filter( filts )
				}
				if *printDockerfile || *printInfo {
					fmt.Println( "Printing dockerfiles" )
					for _,e := range c.Results {
						fmt.Println( "Name: %s\nDescription: %s\nDockerfile:\n%s\n\n", 
							e.Name, e.Description, e.Dockerfile )
					}
				} else {
					fmt.Println( "Search results: ", c.Results )
				}
				
			} else {
				fmt.Println( "No configuration file found, use --generate-config" )
				
			}
		}
	}
}
