package main

import (
	"fmt"
	"github.com/openinfradev/registry-builder/builder/config"
	"github.com/openinfradev/registry-builder/builder/controller"
	"github.com/openinfradev/registry-builder/builder/docs"
	"github.com/openinfradev/registry-builder/builder/network/server"
	"github.com/openinfradev/registry-builder/builder/repository"
	"github.com/openinfradev/registry-builder/builder/service"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

func main() {

	config.LoadConfig()

	defaultinfo := config.GetConfig().Default

	// programatically set swagger info
	docs.SwaggerInfo.Title = "Builder API"
	docs.SwaggerInfo.Description = "This is a sample server for Builder."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = fmt.Sprintf("%s:%s", defaultinfo.Domain, defaultinfo.Port)
	docs.SwaggerInfo.BasePath = "/v1"

	// server ready
	server := server.New()

	// controller ready
	api := controller.New()
	api.RequestMapping(server)

	// redis builder list sync
	registerService := new(service.RegisterService)
	go registerService.Sync()

	// docker login
	dockerService := new(service.DockerService)
	go dockerService.Login()

	// initialize minio port
	registryRepository := new(repository.RegistryRepository)
	if !registryRepository.CreatePortTableIfExists() {
		panic("Failed to create port temporary table for minio")
	}

	// server run
	server.Run(config.GetConfig().Default.Port)
}
