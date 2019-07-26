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

	_ "github.com/go-sql-driver/mysql"
)

func main() {

	// generate swagger documents
	// swag := "../../bin/swag init"
	// _, err := exec.Command("/bin/sh", "-c", swag).Output()
	// if err != nil {
	// 	logger.ERROR("main.go", "Failed to generate swagger documents")
	// } else {
	// 	logger.INFO("main.go", "Generate Swagger Documents")
	// }
	swag := exec.Command("../../bin/swag", "init")
	err := swag.Run()
	if err != nil {
		logger.ERROR("main.go", "Failed to generate swagger documents")
	} else {
		logger.INFO("main.go", "Generate Swagger Documents")
	}

	// flags
	loglevel := flag.Int("loglevel", 0, "Log Level 0:debug 1:info 2:error")

	dbhost := flag.String("dbhost", "", "Database Host Name")
	dbport := flag.String("dbport", "", "Database Port")
	dbuser := flag.String("dbuser", "", "Database User Name")
	dbpass := flag.String("dbpass", "", "Database User Password")
	dbname := flag.String("dbname", "", "Database Name")
	dbxarg := flag.String("dbxarg", "", "Database Extra Arguments")

	registryInsecure := flag.Bool("registry-insecure", false, "Docker Registry Insecure")
	registryEndpoint := flag.String("registry-endpoint", "localhost:5000", "Docker Registry Endpoint")

	port := flag.String("port", "4000", "Builder Service Port")

	flag.Parse()

	logger.DEBUG("main.go", fmt.Sprintf("flags information\n loglevel[%d]\n dbhost[%v]\n dbport[%v]\n dbuser[%v]\n dbpass[%v]\n dbname[%v]\n dbxarg[%v]", *loglevel, *dbhost, *dbport, *dbuser, *dbpass, *dbname, *dbxarg))

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
		RegistryInsecure: *registryInsecure,
		RegistryEndpoint: *registryEndpoint,
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
