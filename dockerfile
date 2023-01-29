FROM golang:1.16-alpine

LABEL version="1.0"
LABEL description="HMRP Forum"
LABEL authors="HMRP"
LABEL author-usernames="HMRP"

# for bash troubleshooting
RUN apk add --no-cache bash \
    # Important: required for go-sqlite3
    gcc \
    # Required for Alpine
    musl-dev

COPY . /go/src/app

WORKDIR /go/src/app
COPY . .

RUN go mod download

# build current directory as main
RUN go build -o /forum

# define the port number the container should expose, same as port exposed in main.go
EXPOSE 8080

# run the go file that was built
CMD [ "/forum" ]