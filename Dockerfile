FROM debian:wheezy

RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates

ENV DEVELOPMENT false

ADD . /kb
WORKDIR /kb/

ENV PORT 80
EXPOSE 80

RUN ["chmod", "+x", "/kb/.bin/run"]
CMD ["/kb/.bin/run"]