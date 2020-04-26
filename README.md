
# challenge

Fetch an image (default is https://upload.wikimedia.org/wikipedia/commons/f/fb/Berliner.Philharmonie.von.Sueden.jpg) and deliver it as fast as possible to the client and save it localy on disc at the same time. The local copy is called output.jpg and is saved in the caller directory.

## Used libraries

* [Gorilla Mux](https://github.com/gorilla/mux) - https://github.com/gorilla/mux
* [RxGo](https://github.com/reactivex/rxgo/v2) - https://github.com/reactivex/rxgo/v2
* [Testify](https://github.com/stretchr/testify) - https://github.com/stretchr/testify
* [Testify](https://github.com/namsral/flag) - https://github.com/namsral/flag

## Usage

You can either run the programm within a docker container or build it manually, then you need go version > 1.13 installed. For both instructions run the following commands inside of the cloned folder. 

The programm provide following options:

```shell
-buffersize reading buffer size in bytes
-timeout    http client timeout in seconds
-imageUrl   overwrite default image url
```

### Manual installation

```shell
go build
./challenge
```

### With Docker

I tested it with docker 19.03.06 on Linux and Mac

```shell
docker build -t challenge .
docker run -p 8080:8080 challenge
```

# Testing

You can rann all tests which are provided at once with:

```shell
go test ./...
```
