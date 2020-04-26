FROM golang:1.14.2-buster as builder

# RUN apt-get update && apt-get install git -y

RUN mkdir /app
WORKDIR /app

ENV GO111MODULE=on

COPY . .

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o challenge

# Run container
FROM debian:buster

RUN apt-get update && apt-get install ca-certificates -y

RUN mkdir /app
WORKDIR /app


COPY testdata/feldberg.jpg /app/testdata/feldberg.jpg
COPY --from=builder /app/challenge .

CMD ["./challenge"]
