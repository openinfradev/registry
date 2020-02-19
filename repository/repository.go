package repository

import (
	"database/sql"
	"fmt"
	"github.com/openinfradev/registry-builder/config"
	"github.com/openinfradev/registry-builder/util/logger"
)

// CreateDBConnection return created database connection
func CreateDBConnection() *sql.DB {
	url := ""

	dbinfo := config.GetConfig().Database

	switch dbinfo.Type {
	case "mysql":
		url = fmt.Sprintf("%v:%v@tcp(%v:%v)/%v", dbinfo.User, dbinfo.Password, dbinfo.Host, dbinfo.Port, dbinfo.Name)
		if dbinfo.Xargs != "" {
			url += "?" + dbinfo.Xargs
		}
	case "postgres":
		url = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable options='%s'", dbinfo.Host, dbinfo.Port, dbinfo.User, dbinfo.Password, dbinfo.Name, dbinfo.Xargs)
	}

	db, err := sql.Open(dbinfo.Type, url)
	if err != nil {
		logger.ERROR("repository/repository.go", "CreateDBConnection", "failed database connection : "+err.Error())
	}
	// logger.DEBUG("repository/repository.go", "CreateDBConnection", "created database connection : "+url)

	return db
}

// CloseDBConnection is closing database connection
func CloseDBConnection(dbconn *sql.DB) {
	// logger.DEBUG("repository/repository.go", "CloseDBConnection", "closed database connection")

	dbconn.Close()
}
