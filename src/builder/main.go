package main

import (
	"builder/controller"
	"builder/docs"
	"builder/network/server"
	"builder/repository"
	"builder/service"
	"builder/util/logger"
	"flag"
	"fmt"
	"os/exec"
	"runtime"

	_ "github.com/go-sql-driver/mysql"
)

func main() {

	// generate swagger documents
	if runtime.GOOS != "windows" {
		swag := exec.Command("../../bin/swag", "init")
		err := swag.Run()
		if err != nil {
			logger.ERROR("main.go", "Failed to generate swagger documents")
		} else {
			logger.INFO("main.go", "Generate Swagger Documents")
		}
	} else {
		logger.INFO("main.go", "Skipped Generate Swagger Documents in Windows")
	}

	// flags
	loglevel := flag.Int("log.level", 0, "Log Level 0:debug 1:info 2:error")

	dbhost := flag.String("db.host", "", "Database Host Name")
	dbport := flag.String("db.port", "", "Database Port")
	dbuser := flag.String("db.user", "", "Database User Name")
	dbpass := flag.String("db.pass", "", "Database User Password")
	dbname := flag.String("db.name", "", "Database Name")
	dbxarg := flag.String("db.xarg", "", "Database Extra Arguments")

	registryName := flag.String("registry.name", "registry", "Docker Registry Container Name")
	registryInsecure := flag.Bool("registry.insecure", false, "Docker Registry Insecure")
	registryEndpoint := flag.String("registry.endpoint", "localhost:5000", "Docker Registry Endpoint")

	port := flag.String("service.port", "4000", "Builder Service Port")
	tmpPath := flag.String("service.tmp", "/tmp/builder", "Builder Service Temporary Path")

	flag.Parse()

	logger.DEBUG("main.go", fmt.Sprintf("settings basic\n log.level[%d]\n service.port[%v]\n service.tmp[%v]", *loglevel, *port, *tmpPath))
	logger.DEBUG("main.go", fmt.Sprintf("settings database\n db.host[%v]\n db.port[%v]\n db.user[%v]\n db.pass[%v]\n db.name[%v]\n db.xarg[%v]", *dbhost, *dbport, *dbuser, *dbpass, *dbname, *dbxarg))
	logger.DEBUG("main.go", fmt.Sprintf("settings registry\n registry.name[%v]\n registry.insecure[%v]\n registry.endpoint[%v]", *registryName, *registryInsecure, *registryEndpoint))

	if *dbhost == "" {
		logger.FATAL("main.go", "Required Database Host Name")
	}

	dbinfo := repository.DBInfo{
		DBhost: *dbhost,
		DBport: *dbport,
		DBuser: *dbuser,
		DBpass: *dbpass,
		DBname: *dbname,
		DBxarg: *dbxarg,
	}

	basicinfo := service.BasicInfo{
		RegistryName:     *registryName,
		RegistryInsecure: *registryInsecure,
		RegistryEndpoint: *registryEndpoint,
		TemporaryPath:    *tmpPath,
	}

	// log level
	logger.SetLevel(*loglevel)

	// programatically set swagger info
	docs.SwaggerInfo.Title = "Builder API"
	docs.SwaggerInfo.Description = "This is a sample server for Builder."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:" + *port
	docs.SwaggerInfo.BasePath = "/v1"

	// server ready
	server := server.New()

	// controller ready
	api := controller.New()
	api.RequestMapping(server)

	// database connection ready
	repository.SetDBConnectionInfo(&dbinfo)

	// service info
	service.SetBasicInfo(&basicinfo)

	// server run
	server.Run(*port)
}
