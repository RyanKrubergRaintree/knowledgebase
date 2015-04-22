FROM golang:1.4

ADD . /go/src/github.com/raintreeinc/knowledgebase
RUN go get github.com/raintreeinc/knowledgebase

ADD . /kb

RUN go build -o /kb/knowledgebase github.com/raintreeinc/knowledgebase

ENV CLIENTDIR /kb/client
WORKDIR /kb/

# expose services
ENV PORT 80
EXPOSE 80

CMD ["/kb/knowledgebase", "-config", "/kb/knowledgebase.toml"]