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
	logger.DEBUG("repository/repository.go", "SetDBConnectionInfo", "setting database connection information")

	dbinfo = info
}

// CreateDBConnection return created database connection
func CreateDBConnection() *sql.DB {
	url := ""
	switch dbinfo.DBtype {
	case "mysql":
		url = fmt.Sprintf("%v:%v@tcp(%v:%v)/%v", dbinfo.DBuser, dbinfo.DBpass, dbinfo.DBhost, dbinfo.DBport, dbinfo.DBname)
		if dbinfo.DBxarg != "" {
			url += "?" + dbinfo.DBxarg
		}
	case "postgres":
		url = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable options='%s'", dbinfo.DBhost, dbinfo.DBport, dbinfo.DBuser, dbinfo.DBpass, dbinfo.DBname, dbinfo.DBxarg)
	}

	db, err := sql.Open(dbinfo.DBtype, url)
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
