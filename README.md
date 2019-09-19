Taco-Registry Builder Project
=============================

## Golang

* golang sdk download & install : https://golang.org/dl/
* export GOROOT, GOPATH
```
$ vi ~/.profile
  # golang path
  export GOROOT=<go sdk path>
  export GOPATH=<project path>
  export PATH=$GOROOT:$GOPATH/bin:$PATH

$ source ~/.profile
$ go env
```

## IDE (Visual Studio Code)

* Visual Studio Code install : https://code.visualstudio.com/download
* Extensions install : keyword is "go"

## Environment

* gcc install (ubuntu amd64)
```
$ sudo apt -y update
$ sudo apt install -y build-essential
```
* project dependency library
```
$ make deps
```

## Binary Build

#### Build
```
$ make build
```

## Binary Deploy

#### Run (Only : Ubuntu 18.04 amd64 arch)
```
$ ./builder \
    -log.level=0 \
    -db.type=postgres \
    -db.host=exntu.kr \
    -db.port=25432 \
    -db.user=registry \
    -db.pass=registry1234\$\$ \
    -db.name=registry \
    -db.xarg= \
    -registry.name=taco-registry \
    -registry.insecure=true \
    -registry.endpoint=exntu.kr:25000 \
    -redis.endpoint=exntu.kr:26379 \
    -service.domain=localhost \
    -service.port=4000 \
    -service.tmp=/tmp
```

## Docker Build

#### Git clone
``` 
$ git clone https://starlkj@tde.sktelecom.com/stash/scm/oreotools/taco-registry-builder.git
```

#### Docker build
```
(move taco-registry-builder directory)
$ make docker-build
```

## Docker Deploy

#### Docker run
```
$ docker run -d -p 4000:4000 \
    -v /var/run/docker.sock:/var/run/docker.sock \
    --restart=always --name builder \
    taco-registry/builder:v2 \
    -log.level=0 \
    -db.type=postgres \
    -db.host=exntu.kr \
    -db.port=25432 \
    -db.user=registry \
    -db.pass=registry1234\$\$ \
    -db.name=registry \
    -db.xarg= \
    -registry.name=taco-registry \
    -registry.insecure=true \
    -registry.endpoint=exntu.kr:25000 \
    -redis.endpoint=exntu.kr:26379 \
    -service.domain=localhost \
    -service.port=4000 \
    -service.tmp=/tmp
```