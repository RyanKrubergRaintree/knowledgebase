FROM golang:1.4

# Mongo
RUN apt-key adv --keyserver hkp://keyserver.ubuntu.com:80 --recv 7F0CEB10
RUN echo 'deb http://downloads-distro.mongodb.org/repo/ubuntu-upstart dist 10gen' | tee /etc/apt/sources.list.d/mongodb.list

RUN apt-get -y update && apt-get install -y mongodb-org supervisor

# knowledgebase
RUN env

RUN mkdir -p /go/src/github.com/raintreeinc/knowledgebase
RUN mkdir -p /var/log

COPY . /go/src/github.com/raintreeinc/knowledgebase

RUN go get github.com/raintreeinc/knowledgebase/...
RUN go install github.com/raintreeinc/knowledgebase

ENV DATABASE="mongodb://localhost/knowledgebase"

# Setup supervisord
RUN mkdir -p /var/log/supervisord
COPY supervisord.conf /etc/supervisor/conf.d/supervisord.conf

# expose services
EXPOSE 80

# start deamon
CMD ["/usr/bin/supervisord"]