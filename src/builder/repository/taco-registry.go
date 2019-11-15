package repository

import (
	"builder/constant/minio"
	"builder/model"
	"builder/util/logger"
	"fmt"
)

// RegistryRepository is registry db (postregs)
type RegistryRepository struct{}

// UpdateBuildPhase is build phase changing
func (a *RegistryRepository) UpdateBuildPhase(buildID string, phase string) bool {
	dbconn := CreateDBConnection()
	defer CloseDBConnection(dbconn)

	_, err := dbconn.Exec("update build set phase=$1 where id=$2", phase, buildID)
	if err != nil {
		logger.ERROR("repository/taco-registry.go", "UpdateBuildPhase", err.Error())
		logger.ERROR("repository/taco-registry.go", "UpdateBuildPhase", buildID+"::"+phase)
		return false
	}
	return true
}

// InsertBuildLog is build log insert row to row
func (a *RegistryRepository) InsertBuildLog(row *model.BuildLogRow) bool {
	dbconn := CreateDBConnection()
	defer CloseDBConnection(dbconn)

	if row.Valid {
		_, err := dbconn.Exec("insert into build_log (build_id, seq, type, message, datetime) values ($1, $2, $3, $4, now())", row.BuildID, row.Seq, row.Type, row.Message)
		if err != nil {
			logger.ERROR("repository/taco-registry.go", "InsertBuildLog", err.Error())
			logger.ERROR("repository/taco-registry.go", "InsertBuildLog", row.BuildID+"::"+row.Message)
			return false
		}
		return true
	}
	return false
}

// InsertBuildLogBatch is build log insert rows batch - wrong
func (a *RegistryRepository) InsertBuildLogBatch(rows []model.BuildLogRow) {
	dbconn := CreateDBConnection()
	defer CloseDBConnection(dbconn)

	// rows count -> failed count??
	for _, row := range rows {
		if row.Valid {
			_, err := dbconn.Exec("insert into build_log (build_id, seq, type, message, datetime) values ($1, $2, $3, $4, now())", row.BuildID, row.Seq, row.Type, row.Message)
			if err != nil {
				logger.ERROR("repository/taco-registry.go", "InsertBuildLogBatch", err.Error())
			}
		}
	}
}

// UpdateTagDigest is digest and size updating in tag table
func (a *RegistryRepository) UpdateTagDigest(buildID string, tag string, digest string, size string) bool {
	dbconn := CreateDBConnection()
	defer CloseDBConnection(dbconn)

	_, err := dbconn.Exec("update tag set manifest_digest=$1, size=$2 where build_id=$3 and name=$4 and (end_time is null or end_time > now())", digest, size, buildID, tag)
	if err != nil {
		logger.ERROR("repository/taco-registry.go", "UpdateTagDigest", err.Error())
		return false
	}
	return true
}

// DeleteUsageLog is usage log deleting
func (a *RegistryRepository) DeleteUsageLog(buildID string) bool {
	dbconn := CreateDBConnection()
	defer CloseDBConnection(dbconn)

	_, err := dbconn.Exec("delete from usage_log where build_id=$1 and kind='create_tag' and tag='latest'", buildID)
	if err != nil {
		logger.ERROR("repository/taco-registry.go", "DeleteUsageLog", err.Error())
		return false
	}
	return true
}

// DeleteTag is tag deleting
func (a *RegistryRepository) DeleteTag(buildID string, tag string) bool {
	dbconn := CreateDBConnection()
	defer CloseDBConnection(dbconn)

	_, err := dbconn.Exec("delete from tag where build_id=$1 and name=$2", buildID, tag)
	if err != nil {
		logger.ERROR("repository/taco-registry.go", "DeleteTag", err.Error())
		return false
	}
	return true
}

// CreatePortTableIfExists is creating table for port if exists
func (a *RegistryRepository) CreatePortTableIfExists() bool {

	dbconn := CreateDBConnection()
	defer CloseDBConnection(dbconn)

	_, err := dbconn.Exec("create table if not exists port_temp (port integer not null, primary key (port) ) ")
	if err != nil {
		logger.ERROR("repository/taco-registry.go", "CreatePortTableIfExists", err.Error())
		return false
	}
	return true

}

// GetTopPort returns top port of exists ports
func (a *RegistryRepository) GetTopPort() int {

	dbconn := CreateDBConnection()
	defer CloseDBConnection(dbconn)

	// select top. if not exists port then returns minimum port
	row := dbconn.QueryRow("select port from port_temp order by port desc limit 1")
	if row == nil {
		logger.DEBUG("repository/taco-registry.go", "GetTopPort", "Port is not exists")
		return minio.MinioMinPort
	}
	var topPort int
	err := row.Scan(&topPort)
	if err != nil {
		logger.ERROR("repository/taco-registry.go", "GetTopPort", err.Error())
		return minio.MinioMinPort
	}

	logger.DEBUG("repository/taco-registry.go", "GetTopPort", fmt.Sprintf("Current Top Port is [%v]", topPort))

	return topPort
}

// InsertPort is inserting temporary port
func (a *RegistryRepository) InsertPort(port int) bool {
	dbconn := CreateDBConnection()
	defer CloseDBConnection(dbconn)

	_, err := dbconn.Exec("insert into port_temp (port) values ($1)", port)
	if err != nil {
		logger.ERROR("repository/taco-registry.go", "InsertPort", err.Error())
		return false
	}

	return false
}

// DeletePort is deleting temporary port
func (a *RegistryRepository) DeletePort(port int) bool {
	dbconn := CreateDBConnection()
	defer CloseDBConnection(dbconn)

	_, err := dbconn.Exec("delete from port_temp where port=$1", port)
	if err != nil {
		logger.ERROR("repository/taco-registry.go", "InsertPort", err.Error())
		return false
	}

	return false
}
