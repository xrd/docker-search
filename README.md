# README #

docker-search: a better way to search the docker registry.


### The Configuration File ###

Why is there a configuration file? Well, the docker registry/index/hub search is a little confusing right now. 
The API documents indicated I should be able to use hub.docker.com to search, but then all these requests are either
301 or 302. Something is wrong. So, I sniffed the traffic from `docker search` to see what they were using (with this command `tcpdump -s 0 -A 'tcp[((tcp[12:1] & 0xf0) >> 2):4] = 0x47455420'`). They were using a strange IP and port, and I got it working with the following cURL command `curl -i 192.168.59.103:2375/v1.12/images/search?term=ffmpeg`. I figured this would change soon, and did not want to hard code the URL into the app. So, docker-search uses a configuration file which you can edit if the backend changes or solidifies.
### Contribution guidelines ###

