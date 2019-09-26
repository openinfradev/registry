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
FROM docker:dind AS image
LABEL maintainer="linus lee <linus@exntu.com>"

WORKDIR /
COPY --from=build /work/builder .

EXPOSE 4000

ENTRYPOINT ["./builder"]
