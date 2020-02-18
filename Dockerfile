# Taco-Registry Build Stage
FROM golang:latest AS build
LABEL maintainer="Seungkyu Ahn <seungkyua@gmail.com>"

RUN mkdir -p /{work,go_workspace}
WORKDIR /work

COPY go.mod .
COPY go.sum .
RUN go mod vendor
COPY . .

ENV GOPATH /go_workspace
RUN make swag
RUN make build


# Taco-Registry Image Stage
FROM ubuntu:18.04 AS image
LABEL maintainer="linus lee <linus@exntu.com>"

RUN sed -i 's/archive.ubuntu.com/ftp.neowiz.com\/ubuntu/g' /etc/apt/sources.list \
    && apt-get -y update \
    && apt-get -y install apt-transport-https ca-certificates curl gnupg-agent software-properties-common \
    && curl -fsSL https://download.docker.com/linux/ubuntu/gpg | apt-key add - \
    && apt-key fingerprint 0EBFCD88 \
    && add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable" \
    && apt-get -y update \
    && apt-get -y install docker-ce docker-ce-cli containerd.io git \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /
COPY --from=build /work/bin/builder .
COPY --from=build /work/builder/docs .

RUN mkdir -p /conf
COPY --from=build /work/builder/conf/* /conf/

EXPOSE 4000

ENTRYPOINT ["./builder"]
