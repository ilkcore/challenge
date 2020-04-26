
# challenge

Fetch an image (default is https://upload.wikimedia.org/wikipedia/commons/f/fb/Berliner.Philharmonie.von.Sueden.jpg) and deliver it as fast as possible to the client and save it localy on disc at the same time. The local copy is called output.jpg and is saved in the caller directory.

## Used libraries

* [Gorilla Mux](https://github.com/gorilla/mux) - https://github.com/gorilla/mux
* [RxGo](https://github.com/reactivex/rxgo/v2) - https://github.com/reactivex/rxgo/v2
* [Testify](https://github.com/stretchr/testify) - https://github.com/stretchr/testify
* [Flag](https://github.com/namsral/flag) - https://github.com/namsral/flag

## Installation

You can either run the programm within a docker container or build it manually, then you need go version > 1.13 installed.

The programm provide following options, either via commandline argument or via ENV variable:

|Cmd arg|ENV|Type|Description|
|:---|:---|:---|:---|
|-buffersize| BUFFERSIZE| int| reading buffer size in bytes|
|-timeout| TIMEOUT| int|http client timeout in seconds|
|-imageurl| IMAGEURL| string|overwrite default image url|

### Manual

```shell
go build
./challenge
```

### With Docker

I tested it with docker 19.03.06 on Linux and Mac

```shell
docker build -t challenge .
docker run -p 8080:8080 challenge

# To run docker in background use
docker run -p 8080:8080 -d challenge

# To provide ENV variable to docker use the following command
docker run -p 8080:8080 -e BUFFERSIZE=1024 -e TIMEOUT=60 challenge

# If you want to receive the image copy from docker container use this command:
# use 'docker ps' to get die container id
docker cp CONTAINER_ID:/app/output.jpg docker-copy.jpg
```

# Usage

Use the browser of your choise and open http://localhost:8080

Or use curl 
```shell
curl -v -o download.jpg http://localhost:8080
```

# Testing

You can run all code tests which are provided at once with in the root of the directory.

```shell
go test ./...
```
