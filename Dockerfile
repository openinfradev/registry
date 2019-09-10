# Taco-Registry Builder
FROM golang:latest

LABEL maintainer="linus lee <linus@exntu.com>"

WORKDIR /go/src/builder

ADD src src
#COPY ./src/builder ./

#RUN go mod download

RUN go build -o builder ./src/builder/main.go

EXPOSE 4000

CMD ["./builder"]
