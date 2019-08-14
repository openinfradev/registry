package repository

import (
	"builder/model"
	"builder/util/logger"
	"fmt"
)

// RegistryRepository is registry db (postregs)
type RegistryRepository struct{}

// SelectCommonCodeList is test
func (a *RegistryRepository) SelectCommonCodeList() []model.RegistryCommonCode {
	dbconn := CreateDBConnection()
	defer CloseDBConnection(dbconn)

	codeName := ""
	groupCode := ""

	rows, err := dbconn.Query("select code_name as codeName, group_code as groupCode from common_code")
	if err != nil {
		logger.ERROR("registry.go", err.Error())
		return []model.RegistryCommonCode{}
	}
	defer rows.Close()

	codeList := []model.RegistryCommonCode{}
	for rows.Next() {
		err := rows.Scan(&codeName, &groupCode)
		if err != nil {
			continue
		}
		code := model.RegistryCommonCode{
			CodeName:  codeName,
			GroupCode: groupCode,
		}
		codeList = append(codeList, code)
	}

	// debug
	for i, c := range codeList {
		logger.DEBUG("registry.go", fmt.Sprintf("read code seq[%d] codeName[%v] groupCode[%v]", i, c.CodeName, c.GroupCode))
	}

	return codeList
}
