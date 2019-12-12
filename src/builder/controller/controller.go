package controller

import (
	"builder/network/server"
	"builder/service"
	"github.com/gin-gonic/gin"
)

// RequestMapper is path and request method set
type RequestMapper struct {
	Method  string
	Path    string
	Request gin.HandlerFunc
}

// RestAPI is request mapping functions
type RestAPI struct {
	RequestMapping func(*server.Instance)
}

// InjectedServices is injection services
type InjectedServices struct {
	DockerService    *service.DockerService
	RegistryService  *service.RegistryService
	MinioService     *service.MinioService
	SecurityService  *service.SecurityService
}

var is *InjectedServices

var mappers []RequestMapper

func init() {
	// inject service
	is = &InjectedServices {
		DockerService:   new(service.DockerService),
		RegistryService: new(service.RegistryService),
		MinioService:    new(service.MinioService),
		SecurityService: new(service.SecurityService),
	}
	
}

// New returns RestAPI
func New() *RestAPI {
	api := &RestAPI{
		RequestMapping: requestMapping,
	}
	return api
}

func requestMapping(instance *server.Instance) {

	for _, mapper := range mappers {
		instance.AddRoute(mapper.Method, mapper.Path, mapper.Request)
		// logger.INFO("controller.go", fmt.Sprintf("[%v] [%v] mapped request method", mapper.Method, mapper.Path))
	}
}

func addRequestMapping(mapper RequestMapper) {
	mappers = append(mappers, mapper)
}

func authorization(token string) bool {
	if token == "" {
		return false
	}
	return service.Authorization(token);
}