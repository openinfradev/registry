package service

import (
	"builder/repository"
	"builder/util/logger"
)

// SampleService is sample service
type SampleService struct{}

var sampleRepository *repository.SampleRepository

func init() {
	// inject repositories
	injectRepositories()
}

func injectRepositories() {
	sampleRepository = new(repository.SampleRepository)
}

// GetWorkflowList returns workflow list by keyword
func (s *SampleService) GetWorkflowList(keyword string) []repository.Workflow {
	// sampleRepository := new(repository.SampleRepository)
	return sampleRepository.GetWorkflowList(keyword)
}

// Holding is sleep 3 seconds
func (s *SampleService) Holding(target string) string {
	// sampleRepository := new(repository.SampleRepository)
	logger.DEBUG("service/sample.go", target+" sleep seconds started")
	r := sampleRepository.Holding(target)
	logger.DEBUG("service/sample.go", r)
	return r
}
