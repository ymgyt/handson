FROM golang:alpine3.7

# see https://github.com/docker-library/golang/blob/1e2ddb8ec8c9f59a56ddfb343e1bd8e65440b6db/1.9/alpine3.7/Dockerfile#L54-L58
ADD . /go/src/github.com/YmgchiYt/golang-handson/web_development_with_go

WORKDIR /go/src/github.com/YmgchiYt/golang-handson/web_development_with_go
RUN apk update && apk add git
RUN go get github.com/golang/dep/cmd/dep
RUN dep ensure
RUN go install github.com/YmgchiYt/golang-handson/web_development_with_go

ENTRYPOINT /go/bin/web_development_with_go

EXPOSE 3000
