package main

import (
	"builder/controller"
	"builder/network/server"
	"builder/repository"
	"builder/util/logger"
	"flag"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func main() {

	// flags
	loglevel := flag.Int("loglevel", 0, "Log Level 0:debug 1:info 2:error")

	dbhost := flag.String("dbhost", "", "Database Host Name")
	dbport := flag.String("dbport", "", "Database Port")
	dbuser := flag.String("dbuser", "", "Database User Name")
	dbpass := flag.String("dbpass", "", "Database User Password")
	dbname := flag.String("dbname", "", "Database Name")
	dbxarg := flag.String("dbxarg", "", "Database Extra Arguments")

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

	// log level
	logger.SetLevel(*loglevel)

	// server ready
	server := server.New()

	// controller ready
	api := controller.New()
	api.RequestMapping(server)

	// database connection ready
	repository.SetDBConnectionInfo(&dbinfo)

	// server run
	server.Run("4000")
}
