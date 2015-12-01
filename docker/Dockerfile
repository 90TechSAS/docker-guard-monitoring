FROM     ubuntu:trusty
MAINTAINER Nurza nurza@90tech.fr

RUN (apt-get update && apt-get install -y -q wget git curl && apt-get -y -q autoclean && apt-get -y -q autoremove)

RUN (wget -O /tmp/go.tar.gz https://storage.googleapis.com/golang/go1.5.linux-amd64.tar.gz)
RUN (cd /tmp && tar xf go.tar.gz && mv go /usr/local)

RUN mkdir /go
RUN GOPATH=/go /usr/local/go/bin/go get "github.com/influxdb/influxdb/client"
RUN GOPATH=/go /usr/local/go/bin/go get "github.com/nurza/logo"
RUN GOPATH=/go /usr/local/go/bin/go get "github.com/90TechSAS/libgo-docker-guard"
RUN GOPATH=/go /usr/local/go/bin/go get "github.com/gorilla/mux"
RUN GOPATH=/go /usr/local/go/bin/go get "gopkg.in/yaml.v2"

ADD dgm /dgm
RUN (cd /dgm && GOPATH=/go /usr/local/go/bin/go build -o dg-monitoring main.go)

EXPOSE 8124

ENTRYPOINT /dgm/dg-monitoring -vvv -f /dgm/config.yaml
