package repository

import (
	"builder/util/logger"
	"database/sql"
	"fmt"
)

// DBInfo is basically database information
type DBInfo struct {
	DBtype string
	DBhost string
	DBport string
	DBuser string
	DBpass string
	DBname string
	DBxarg string
}

var dbinfo *DBInfo

// SetDBConnectionInfo is setting database basically information
func SetDBConnectionInfo(info *DBInfo) {
	logger.DEBUG("repository.go", "setting database connection information")

	dbinfo = info
}

// CreateDBConnection return created database connection
func CreateDBConnection() *sql.DB {
	url := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v", dbinfo.DBuser, dbinfo.DBpass, dbinfo.DBhost, dbinfo.DBport, dbinfo.DBname)
	if dbinfo.DBxarg != "" {
		url += "?" + dbinfo.DBxarg
	}
	db, err := sql.Open(dbinfo.DBtype, url)
	if err != nil {
		logger.FATAL("repository.go", "failed database connection")
	}
	logger.DEBUG("repository.go", "created database connection")

	return db
}

// CloseDBConnection is closing database connection
func CloseDBConnection(dbconn *sql.DB) {
	logger.DEBUG("repository.go", "closed database connection")

	dbconn.Close()
}
