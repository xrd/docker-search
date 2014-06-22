package main

import (
	"fmt"
	"flag"
	"os"
	"os/user"
	"path"
	"encoding/json"
	"strconv"
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
-format=table      # Format the output, as a table (default) or json
-annotate=true     # Skip annotation (faster, but ignores filters)
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

func generateDefaultConfiguration() bool {
	configPath := getConfigFilePath()
	rv := false
	// Generate configuration
	if file, err := os.Create( configPath ); nil == err {
		def := `
Host = "https://index.docker.io"
Endpoint = "/v1/search"
UpdateCheck = true
`
		file.Write( []byte(def) )
		rv = true
	}
	return rv
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

var genCon *bool
var printDockerfile *bool
// var printInfo *bool
var format *string
var annotation *bool
var filts filters

func main() {

	genCon = flag.Bool( "generate-config", false, "Generate a new default configuration file" )
	printDockerfile = flag.Bool( "dockerfile", false, "Print out dockerfiles" )
	// printInfo = flag.Bool( "info", false, "Print out detailed information on the maintainer(s)" )
	format = flag.String( "format", "table", "Format the output: table or json" )
	annotation = flag.Bool( "annotate", true, "Annotation with Dockerfile information (faster without but no second level search)" )
	flag.Var( &filts, "filter", "List of filters" )
	flag.Parse()

	if *genCon {
		if generateDefaultConfiguration() {
			fmt.Println( "Generated configuration file." )
		} else  {
			fmt.Println( "Unable to create configuration file." )
		}
		// help()
		// flag.PrintDefaults()
	} else {
		if "" == flag.Arg(0) {
			help()
		} else {
			c := new(Client)
			if c.LoadConfig( getConfigFilePath() ) {
				count := 0
				ch := make(chan bool)
				for _, q := range flag.Args() {
					go query( c, q, ch ) //  c.Query(
					//  flag.Arg(0) )
					count++
				}

				if *annotation {
					for count > 0 { 
						<- ch
						count--
					}
					c.Annotate()
					c.Filter( filts )
				}

				printResults( c )
				
			} else {
				fmt.Println( "No configuration file found, use --generate-config" )
				
			}
		}
	}
}

func query( c* Client, query string, ch chan bool ) {
	ch <- c.Query( query )
}

const tableFmt = "%-30s%-30s"

func formatTable( c* Client ) {
	count := strconv.Itoa( len(c.Results) )
	fmt.Println( "Found " + count + " results\n\n" )

	for _,e := range c.Results {
		fmt.Println( fmt.Sprintf( tableFmt, "Name: ",  e.Name ) )
		fmt.Println( fmt.Sprintf( tableFmt, "Description: ", e.Description ) )
		fmt.Println( fmt.Sprintf( "Dockerfile\n\n%s\n\n", e.Dockerfile ) )
	}
}

func formatJson( c* Client ) {
	b, _ := json.Marshal( c.Results )
	fmt.Println( string(b) )
}

func printResults( c* Client ) {
	if *printDockerfile { // || *printInfo {
		switch *format {
		case "json": formatJson( c )
		default: formatTable( c )
		}
	} else {
		fmt.Println( fmt.Sprintf(         tableFmt, "Name", "Description" ) )
		fmt.Println( fmt.Sprintf(         tableFmt, "----", "-----------" ) )
		for _, r := range c.Results {
			fmt.Println( fmt.Sprintf( tableFmt, r.Name, r.Description ) )
		}

	}
}
