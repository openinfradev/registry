
deps:
	@echo "### Pulling dependency libraries."
	@echo "================================="
	go get -u github.com/go-sql-driver/mysql
	go get -u github.com/lib/pq
	go get -u github.com/gin-gonic/gin
	go get -u github.com/go-redis/redis
	go get -u github.com/swaggo/swag/cmd/swag
	go get -u github.com/gin-contrib/location
	go get -u github.com/swaggo/gin-swagger
	go get -u github.com/swaggo/gin-swagger/swaggerFiles
	go get -u github.com/alecthomas/template
	go get -u gopkg.in/yaml.v2

swag:
	@echo "### Generating taco-registry Builder Swagger Docs."
	go get -u github.com/swaggo/swag/cmd/swag
	${GOPATH}/bin/swag init

build:
	@echo "### Building taco-registry Builder."
	@echo "==================================="
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a --ldflags=--s -o ./bin/builder main.go

build-cross:
	@echo "### Cross-Compiling taco-registry Builder."
	@echo "Darwin amd64(OSX), Windows amd64, Linux arm64, Linux amd64"
	@echo "=========================================================="
	./cross-compile.sh

docker-build:
	@echo "### Making Builder docker image. Multi-Stage"
	@echo "============================================"
	docker build --network=host --no-cache -t taco/registry-builder:latest . -f ./Dockerfile

docker-build-single:
	@echo "### Making Builder docker image."
	@echo "================================"
	docker build --network=host --no-cache -t taco/registry-builder:latest . -f ./Dockerfile.single