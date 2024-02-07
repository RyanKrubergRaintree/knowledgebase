FROM alpine:3.1

RUN apk add --update ca-certificates && rm -rf /var/cache/apk/*

ENV DEVELOPMENT   false
ENV REDIRECTHTTPS true

ADD . /kb
WORKDIR /kb/

ENV PORT 80
EXPOSE 80

# TODO: remove \/
RUN ["ls", "-a", "/kb"] 
RUN ["chmod", "+x", "/kb/.bin/run"]
# TODO: remove \/
RUN ["ls", "-l", "/kb/.bin/run"]
CMD ["/kb/.bin/run"]