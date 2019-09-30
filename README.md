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

## Swagger Documents
* generate documents
   ```
   $ cd src/builder
   $ ../../bin/swag init
   ```
* swagger documents page url
   ```
   http://<builder domain>:<builder port>/swagger/index.html
   ```

## Binary Build

#### Build
```
$ make build
```

## Binary Deploy

#### Run (Only : Ubuntu 18.04 amd64 arch)
```
$ ./builder
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
    taco-registry/builder:v2 
```

## Environment Variables

> optional value : override configuration

| variable | default value | description |
| ------ | ------ | ------ |
| BUILDER_CONFIG | conf/config.yml | configuration file path |
| BUILDER_LOG_LEVEL | DEBUG | log level : DEBUG, INFO, ERROR(default) | 

## Configuration

#### configuration yaml
```
<project dir>/src/builder/conf/config.yml
```

#### configuration yaml fields
| section | field | type | description |
| ------ | ------ | ------ | ------ |
| default | domain | string | service domain |
| default | port | string | service port |
| default | loglevel | int | 0 is DEBUG, 1 is INFO, 2 is ERROR |
| default | tmp | string | temporary directory |
| database | type | string | database type : mysql, postgres |
| database | host | string | database hostname(domain) |
| database | port | string | database port |
| database | name | string | database name |
| database | user | string | database user ID |
| database | password | string | database user password |
| database | xargs | string | database connection extra arguments |
| registry | name | string | registry name |
| registry | insecure | bool | insecure flag |
| registry | endpoint | string | registry endpoint |
| registry | auth | string | registry authorization token api url |
| redis | endpoint | string | redis endpoint |
| clair | endpoint | string | clair endpoint |

#### configuration example
```
#####################################
# Taco-registry builder configuration
#####################################

# default configuration
default:
  domain: 192.168.201.2
  port: 4000
  # DEBUG, INFO, ERROR
  loglevel: DEBUG
  tmp: "/tmp"

# database configuration
database:
  # postgres, mysql
  type: postgres
  host: exntu.kr
  port: 25432
  name: registry
  user: registry
  password: "registry1234$$"
  xargs:

# docker private registry configuration
registry:
  name: taco-registry
  insecure: true
  endpoint: "exntu.kr:25000"
  auth: "http://exntu.kr:38383/api/oauth/token"

# redis configuration
redis:
  endpoint: "exntu.kr:26379"

# clair configuration
clair:
  endpoint: "exntu.kr:26060"
```