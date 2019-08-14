package config

import (
	"builder/repository"
	"builder/service"
	"builder/util/logger"
	"flag"
	"fmt"
)

// ParseFlags is all flags parsing and returns basicinfo & dbinfo
func ParseFlags() (*service.BasicInfo, *repository.DBInfo) {
	// flags
	loglevel := flag.Int("log.level", 0, "Log Level 0:debug 1:info 2:error")

	dbtype := flag.String("db.type", "mysql", "Database Type (mysql, postgres)")
	dbhost := flag.String("db.host", "", "Database Host Name")
	dbport := flag.String("db.port", "", "Database Port")
	dbuser := flag.String("db.user", "", "Database User Name")
	dbpass := flag.String("db.pass", "", "Database User Password")
	dbname := flag.String("db.name", "", "Database Name")
	dbxarg := flag.String("db.xarg", "", "Database Extra Arguments")

	registryName := flag.String("registry.name", "registry", "Docker Registry Container Name")
	registryInsecure := flag.Bool("registry.insecure", false, "Docker Registry Insecure")
	registryEndpoint := flag.String("registry.endpoint", "localhost:5000", "Docker Registry Endpoint")

	redisEndpoint := flag.String("redis.endpoint", "localhost:6379", "Redis Endpoint")

	port := flag.String("service.port", "4000", "Builder Service Port")
	tmpPath := flag.String("service.tmp", "/tmp/builder", "Builder Service Temporary Path")

	flag.Parse()

	logger.DEBUG("config.go", fmt.Sprintf("settings basic\n log.level[%d]\n service.port[%v]\n service.tmp[%v]", *loglevel, *port, *tmpPath))
	logger.DEBUG("config.go", fmt.Sprintf("settings database\n db.host[%v]\n db.port[%v]\n db.user[%v]\n db.pass[%v]\n db.name[%v]\n db.xarg[%v]", *dbhost, *dbport, *dbuser, *dbpass, *dbname, *dbxarg))
	logger.DEBUG("config.go", fmt.Sprintf("settings registry\n registry.name[%v]\n registry.insecure[%v]\n registry.endpoint[%v]", *registryName, *registryInsecure, *registryEndpoint))
	logger.DEBUG("config.go", fmt.Sprintf("settings redis\n redis.endpoint[%v]", *redisEndpoint))

	if *dbhost == "" {
		logger.FATAL("config.go", "Required Database Host Name")
	}

	// log level
	logger.SetLevel(*loglevel)

	dbinfo := repository.DBInfo{
		DBtype: *dbtype,
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
		RedisEndpoint:    *redisEndpoint,
		ServicePort:      *port,
	}
	return &basicinfo, &dbinfo
}
