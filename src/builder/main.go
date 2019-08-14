package main

import (
	"builder/config"
	"builder/controller"
	"builder/docs"
	"builder/network/server"
	"builder/repository"
	"builder/service"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

func main() {

	// generate swagger documents
	// if runtime.GOOS != "windows" {
	// 	swag := exec.Command("../../bin/swag", "init")
	// 	err := swag.Run()
	// 	if err != nil {
	// 		logger.ERROR("main.go", "Failed to generate swagger documents")
	// 	} else {
	// 		logger.INFO("main.go", "Generate Swagger Documents")
	// 	}
	// } else {
	// 	logger.INFO("main.go", "Skipped Generate Swagger Documents in Windows")
	// }

	basicinfo, dbinfo := config.ParseFlags()

	// programatically set swagger info
	docs.SwaggerInfo.Title = "Builder API"
	docs.SwaggerInfo.Description = "This is a sample server for Builder."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:" + basicinfo.ServicePort
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

	// server run
	server.Run(basicinfo.ServicePort)
}
