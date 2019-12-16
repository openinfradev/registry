# Taco-Registry Build Stage
FROM golang:latest AS build
LABEL maintainer="linus lee <linus@exntu.com>"

RUN mkdir -p /work
ENV GOPATH /work
WORKDIR /work

COPY . .

RUN make deps
RUN make build


# Taco-Registry Image Stage
FROM ubuntu:18.04 AS image
LABEL maintainer="linus lee <linus@exntu.com>"

RUN apt-get -y update
RUN apt-get -y install apt-transport-https ca-certificates curl gnupg-agent software-properties-common
RUN curl -fsSL https://download.docker.com/linux/ubuntu/gpg | apt-key add -
RUN apt-key fingerprint 0EBFCD88
RUN add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable"
RUN apt-get -y update
RUN apt-get -y install docker-ce docker-ce-cli containerd.io git

WORKDIR /
COPY --from=build /work/builder .

RUN mkdir /conf
COPY --from=build /work/src/builder/conf/* /conf/

EXPOSE 4000

ENTRYPOINT ["./builder"]
