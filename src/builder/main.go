package main

import (
	"builder/config"
	"builder/controller"
	"builder/docs"
	"builder/network/server"
	"builder/repository"
	"builder/service"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

func main() {

	basicinfo, dbinfo := config.LoadConfig()

	// programatically set swagger info
	docs.SwaggerInfo.Title = "Builder API"
	docs.SwaggerInfo.Description = "This is a sample server for Builder."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = fmt.Sprintf("%s:%s", basicinfo.ServiceDomain, basicinfo.ServicePort)
	docs.SwaggerInfo.BasePath = "/v1"

	// server ready
	server := server.New()

	// controller ready
	api := controller.New()
	api.RequestMapping(server)

	// database connection ready
	repository.SetDBConnectionInfo(dbinfo)

	// service info
	service.SetBasicInfo(basicinfo)

	// redis builder list sync
	registerService := new(service.RegisterService)
	go registerService.Sync()

	// docker login
	dockerService := new(service.DockerService)
	go dockerService.Login()

	// server run
	server.Run(basicinfo.ServicePort)
}
