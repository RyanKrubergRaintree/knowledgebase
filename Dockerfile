FROM golang:1.4

ADD . /go/src/github.com/raintreeinc/knowledgebase
RUN go get github.com/raintreeinc/knowledgebase

ADD . /knowledgebase

RUN go build -o /knowledgebase/kb github.com/raintreeinc/knowledgebase

ENV CLIENTDIR /knowledgebase/client
WORKDIR /knowledgebase/

# expose services
ENV PORT 80
EXPOSE 80

CMD ["/knowledgebase/kb", "-config", "/knowledgebase/knowledgebase.toml"]