package repository

import (
	"builder/util/logger"
	"fmt"
	"strconv"
	"time"
)

// SampleRepository is sample
type SampleRepository struct{}

// Workflow is id and name of workflow
type Workflow struct {
	WorkflowID   string `json:"workflowId"`
	WorkflowName string `json:"workflowName"`
}

// GetWorkflowList returns workflow list for sample
func (s *SampleRepository) GetWorkflowList(keyword string) []Workflow {
	dbconn := CreateDBConnection()
	defer CloseDBConnection(dbconn)

	id := ""
	name := ""

	rows, err := dbconn.Query("select workflow_id as id, workflow_name as name from T_WORKFLOW where workflow_name like concat('%',?,'%')", keyword)
	if err != nil {
		return []Workflow{}
	}
	defer rows.Close()

	workflowList := []Workflow{}
	for rows.Next() {
		err := rows.Scan(&id, &name)
		if err != nil {
			continue
		}
		workflow := Workflow{
			WorkflowID:   id,
			WorkflowName: name,
		}
		workflowList = append(workflowList, workflow)
	}

	// debug
	for idx, workflow := range workflowList {
		logger.DEBUG("sample.go", fmt.Sprintf("read workflow seq[%d] id[%v] name[%v]", idx, workflow.WorkflowName, workflow.WorkflowID))
	}

	return workflowList
}

// Holding is sleep 3 seconds
func (s *SampleRepository) Holding(target string) string {
	logger.DEBUG("repository/sample.go", target+" sleep seconds started")
	i, _ := strconv.Atoi(target)
	time.Sleep(time.Second * time.Duration(i))
	logger.DEBUG("repository/sample.go", target+" sleep seconds ended")
	return target + " sleep seconds ended"
}
