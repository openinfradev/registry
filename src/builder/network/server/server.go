package server

import (
	"builder/util/logger"

	"github.com/gin-gonic/gin"
)

// Instance is server instance functions
type Instance struct {
	Run      func(string)
	AddRoute func(string, string, gin.HandlerFunc)
}

var r *gin.Engine
var v1 *gin.RouterGroup

// New returns created instance
func New() *Instance {
	r = gin.Default()
	v1 = r.Group("/v1")

	instance := &Instance{
		Run:      run,
		AddRoute: addRoute,
	}
	return instance
}

func run(port string) {
	logger.INFO("server.go", "started server on :"+port)
	r.Run(":" + port)
}

func addRoute(method string, path string, routerFunc gin.HandlerFunc) {
	switch method {
	case "GET":
		v1.GET(path, routerFunc)
	case "POST":
		v1.POST(path, routerFunc)
	case "PUT":
		v1.PUT(path, routerFunc)
	case "DELETE":
		v1.DELETE(path, routerFunc)
	default:
	}
}
