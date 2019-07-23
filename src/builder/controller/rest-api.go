package controller

import (
	"builder/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

var sampleService *service.SampleService
var httpsampleService *service.HTTPSampleService

func init() {
	// inject service
	injectServices()

	// ping
	addRequestMapping(
		RequestMapper{
			Method:  "GET",
			Path:    "/health",
			Request: health,
		},
	)
	// sample
	addRequestMapping(
		RequestMapper{
			Method:  "GET",
			Path:    "/sample",
			Request: sample,
		},
	)
	// sleepTest
	addRequestMapping(
		RequestMapper{
			Method:  "GET",
			Path:    "/sleep",
			Request: sleepTest,
		},
	)
	// httpTest
	addRequestMapping(
		RequestMapper{
			Method:  "GET",
			Path:    "/http",
			Request: httpTest,
		},
	)
}

func injectServices() {
	sampleService = new(service.SampleService)
	httpsampleService = new(service.HTTPSampleService)
}

/*
	Request Mapping Functions
*/

func health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
}

func sample(c *gin.Context) {
	keyword := c.Query("q")

	workflowList := sampleService.GetWorkflowList(keyword)

	c.JSON(http.StatusOK, workflowList)
}

func sleepTest(c *gin.Context) {
	target := c.Query("t")

	r := sampleService.Holding(target)

	c.JSON(http.StatusOK, gin.H{
		"message": r,
	})
}

func httpTest(c *gin.Context) {
	r := httpsampleService.GetDaum()

	c.Data(http.StatusOK, "text/html; charset=utf-8", r)
}

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
