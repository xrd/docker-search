# README #

docker-search: a better way to search the docker registry.

### Why? 

This does not seem optimal:

    $ docker search libavcodec
    NAME      DESCRIPTION   STARS     OFFICIAL   AUTOMATED
    
    $ docker search ffmpeg
    NAME                            DESCRIPTION                                     STARS     OFFICIAL   AUTOMATED
    cellofellow/ffmpeg                                                              1                    [OK]
    bfirsh/ffmpeg                                                                   0                    [OK]
    robd/aws-ffmpeg                                                                 0                    [OK]
    lmars/ffmpeg                    FFMpeg built from source (git://source.ffm...   0                    
    link/ffmpeg-built               this has python devenv, a few other build ...   0                    
    miovision/ffmpeg                                                                0                    
    paulbrennan/ffmpeg                                                              0                    
    cmark/ubuntu-ffmpeg                                                             0                    
    mikehearn/ubuntu-ffmpeg                                                         0                    
    cmark/ubuntu-ffmpeg-ssh                                                         0                    
    cmark/ubuntu-14.04-ffmpeg-nfs                                                   0                    
    asachs/docker-ffmpeg            Docker based FFMPEG - keeping the dependan...   0                    

Which of these have libavcodec? Do I have to manually go to hub.registry.docker.com and search through Dockerfiles?

How about this instead?

    $ docker-search -filter=libavcodec ffmpeg
    Name                          Description                   
    ----                          -----------                   
    bfirsh/ffmpeg 

More about what is happening:

    $ docker-search -filter=libavcodec -verbose=true -dockerfile ffmpeg | grep -C 5 avcodec
    Query docker for ffmpeg
    Query finished for: ffmpeg
    Annotating dockerfile for cellofellow/ffmpeg
    Annotating dockerfile for bfirsh/ffmpeg
    Annotating dockerfile for robd/aws-ffmpeg
    Annotating dockerfile for lmars/ffmpeg
    ...
    Finished annotation of dockerfiles
    Filters:  libavcodec
    Filtering dockerfiles
    Found match inside Dockerfile
    Adding result to results: bfirsh/ffmpeg
    ...
    FROM ubuntu:12.10
    MAINTAINER Ben Firshman "ben@orchardup.com"
    RUN echo "deb http://archive.ubuntu.com/ubuntu quantal main universe" > /etc/apt/sources.list
    RUN apt-get update
    RUN apt-get -y install ffmpeg libavcodec-extra-53
    ...
    

### Usage


    docker-search: A better way to search the docker registry
    
    docker-search does a search against the Docker registry, and then pulls Dockerfiles for 
    matching images, searching inside them for more details.

    Examples:
    
    docker-search -filter=quantal ffmpeg  # Search for the string quantal in the Dockerfile
    docker-search -dockerfile ffmpeg # print out full dockerfiles
    docker-search -dockerfile -format=json ffmpeg # print out full dockerfiles as JSON
    
    Flags:
    
      -annotate=true: Annotation with Dockerfile information (faster without but no second level search)
      -dockerfile=false: Print out dockerfiles
      -filter=: List of filters
      -format="table": Format the output: table or json
      -generate-config=false: Generate a new default configuration file
      -verbose=false: Output verbose messages (false)

### Installation

    go get github.com/xrd/docker-search

Or, download the source and install yourself.

### Developer Details ###

docker-search works in three stages right now inside client.go

* Query: query the docker index
* Annotate: annotate the search result with scraped data, like Dockerfiles themselves
* Filter: filter through the data based on interesting information

This structure allows us to write tests since we can seed data (using the `Load` method 
on the client) that has us skip the first two stages. As long as our results are normal.

### The Configuration File ###

Why is there a configuration file? Well, the docker registry/index/hub search is a little confusing right now. 
The API documents indicated I should be able to use hub.docker.com to search, but then all these requests are either
301 or 302. Something is wrong. So, I sniffed the traffic from `docker search` to see what they were using (with this command `tcpdump -s 0 -A 'tcp[((tcp[12:1] & 0xf0) >> 2):4] = 0x47455420'`). They were using a strange IP and port, and I got it working with the following cURL command `curl -i 192.168.59.103:2375/v1.12/images/search?term=ffmpeg`. I figured this would change soon, and did not want to hard code the URL into the app. So, docker-search uses a configuration file which you can edit if the backend changes or solidifies.
### Contribution guidelines ###

