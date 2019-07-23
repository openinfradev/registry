package controller

import (
	"builder/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// var sampleService *service.SampleService
// var httpsampleService *service.HTTPSampleService
var dockerService *service.DockerService

func init() {
	// inject service
	injectServices()

	// health
	addRequestMapping(
		RequestMapper{
			Method:  "GET",
			Path:    "/health",
			Request: health,
		},
	)

	// docker catalog
	addRequestMapping(
		RequestMapper{
			Method:  "GET",
			Path:    "/docker/catalog",
			Request: getDockerCatalog,
		},
	)

}

func injectServices() {
	// sampleService = new(service.SampleService)
	// httpsampleService = new(service.HTTPSampleService)

	dockerService = new(service.DockerService)
}

/*
	Request Mapping Functions
*/

func health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
}

func getDockerCatalog(c *gin.Context) {
	r := dockerService.GetCatalog()

	c.Data(http.StatusOK, "application/json; charset=utf-8", r)
}

// func sample(c *gin.Context) {
// 	keyword := c.Query("q")

// 	workflowList := sampleService.GetWorkflowList(keyword)

// 	c.JSON(http.StatusOK, workflowList)
// }

// func sleepTest(c *gin.Context) {
// 	target := c.Query("t")

// 	r := sampleService.Holding(target)

// 	c.JSON(http.StatusOK, gin.H{
// 		"message": r,
// 	})
// }

// func httpTest(c *gin.Context) {
// 	r := httpsampleService.GetDaum()

// 	c.Data(http.StatusOK, "text/html; charset=utf-8", r)
// }

/*
	registry
*/

/*
	docker build
*/

/*
	security scan
*/

/*
	zookeeper
*/

/*
	docker accout ??
*/
