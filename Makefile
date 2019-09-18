
deps:
	go get -u github.com/go-sql-driver/mysql
	go get -u github.com/lib/pq
	go get -u github.com/gin-gonic/gin
	go get -u github.com/go-redis/redis
	go get -u github.com/swaggo/swag/cmd/swag
	go get -u github.com/gin-contrib/location
	go get -u github.com/swaggo/gin-swagger
	go get -u github.com/swaggo/gin-swagger/swaggerFiles
	go get -u github.com/alecthomas/template

build:
	CGO_ENABLED=0 go build -v --ldflags=--s -o builder src/builder/main.go

docker-build:
	docker build --network host -t taco/builder:latest .
