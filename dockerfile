FROM golang:1.17-alpine as build

LABEL version="1.0"
LABEL description="HMRP Forum"
LABEL authors="HMRP"
LABEL author-usernames="HMRP"

RUN mkdir /app

ADD . /app
# for bash troubleshooting


WORKDIR /app

# for bash troubleshooting
RUN apk add --no-cache bash \
    # Important: required for go-sqlite3
    gcc \
    # Required for Alpine
    musl-dev

RUN go mod download

# build current directory as main
RUN go build -o main .


FROM alpine:3.12

COPY --from=build /app /app
# define the port number the container should expose, same as port exposed in main.gos

WORKDIR /app

# run the go file that was built
CMD ["/app/main"]