package repository

import (
	"builder/model"
	"builder/util/logger"
)

// RegistryRepository is registry db (postregs)
type RegistryRepository struct{}

// SelectCommonCodeList is test
// func (a *RegistryRepository) SelectCommonCodeList() []model.RegistryCommonCode {
// 	dbconn := CreateDBConnection()
// 	defer CloseDBConnection(dbconn)

// 	codeName := ""
// 	groupCode := ""

// 	rows, err := dbconn.Query("select code_name as codeName, group_code as groupCode from common_code")
// 	if err != nil {
// 		logger.ERROR("repository/taco-registry.go", "SelectCommonCodeList", err.Error())
// 		return []model.RegistryCommonCode{}
// 	}
// 	defer rows.Close()

// 	codeList := []model.RegistryCommonCode{}
// 	for rows.Next() {
// 		err := rows.Scan(&codeName, &groupCode)
// 		if err != nil {
// 			continue
// 		}
// 		code := model.RegistryCommonCode{
// 			CodeName:  codeName,
// 			GroupCode: groupCode,
// 		}
// 		codeList = append(codeList, code)
// 	}

// 	// debug
// 	for i, c := range codeList {
// 		logger.DEBUG("repository/taco-registry.go", "SelectCommonCodeList", fmt.Sprintf("read code seq[%d] codeName[%v] groupCode[%v]", i, c.CodeName, c.GroupCode))
// 	}

// 	return codeList
// }

// InsertBuildLog is build log insert row to row
func (a *RegistryRepository) InsertBuildLog(row *model.BuildLogRow) bool {
	dbconn := CreateDBConnection()
	defer CloseDBConnection(dbconn)

	_, err := dbconn.Exec("insert into build_log (build_id, seq, type, message, datetime) values ($1, $2, $3, $4, now())", row.BuildID, row.Seq, row.Type, row.Message)
	if err != nil {
		logger.ERROR("repository/taco-registry.go", "InsertBuildLog", err.Error())
		return false
	}

	return true
}

// InsertBuildLogBatch is build log insert rows batch
func (a *RegistryRepository) InsertBuildLogBatch(rows []model.BuildLogRow) {
	dbconn := CreateDBConnection()
	defer CloseDBConnection(dbconn)

	// rows count -> failed count??
	for _, row := range rows {
		_, err := dbconn.Exec("insert into build_log (build_id, seq, type, message, datetime) values ($1, $2, $3, $4, now())", row.BuildID, row.Seq, row.Type, row.Message)
		if err != nil {
			logger.ERROR("repository/taco-registry.go", "InsertBuildLog", err.Error())
		}
	}

}
