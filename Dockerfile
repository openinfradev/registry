# Taco-Registry Builder
FROM golang:latest
LABEL maintainer="linus lee <linus@exntu.com>"

RUN mkdir -p /work
ENV GOPATH /work
WORKDIR /work

COPY . .

RUN make deps
RUN make build

RUN cp ./builder /
WORKDIR /
RUN rm -rf /work

EXPOSE 4000

CMD ["./builder"]
