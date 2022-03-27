FROM golang:1.18-alpine

## The latest alpine images don't have some tools like (`git` and `bash`).
## Adding git, bash and openssh to the image
RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh

ADD . /src

WORKDIR /src

COPY go.mod go.sum ./

RUN go get -d -v ./...
RUN go install -v ./...

# Build the Go app
RUN go build -o main .

# Run the executable
CMD ["./main"]



## Dockerfile References: https://docs.docker.com/engine/reference/builder/
#
## Start from golang:1.12-alpine base image
#FROM golang:1.18-alpine
#
## The latest alpine images don't have some tools like (`git` and `bash`).
## Adding git, bash and openssh to the image
#RUN #apk update && apk upgrade && \
##    apk add --no-cache bash git openssh
#
## Set the Current Working Directory inside the container
#WORKDIR /src
#
## Copy go mod and sum files
#COPY go.mod go.sum ./
#
## Download all dependancies. Dependencies will be cached if the go.mod and go.sum files are not changed
#RUN go mod download
#
## Copy the source from the current directory to the Working Directory inside the container
#COPY . .
#
## Build the Go app
#RUN go build -o main .
#
## Expose port 8080 to the outside world
##EXPOSE 8080
#
## Run the executable
#CMD ["./main"]