package controller

import (
	"builder/network/server"
	"builder/util/logger"
	"fmt"

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

var mappers []RequestMapper

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
		logger.INFO("controller.go", fmt.Sprintf("[%v] [%v] mapped request method", mapper.Method, mapper.Path))
	}
}

func addRequestMapping(mapper RequestMapper) {
	mappers = append(mappers, mapper)
}
