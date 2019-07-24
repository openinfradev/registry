package controller

import (
	"builder/service"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// var sampleService *service.SampleService
// var httpsampleService *service.HTTPSampleService
var dockerRegistryService *service.DockerRegistryService

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

	// docker repositories tag list
	addRequestMapping(
		RequestMapper{
			Method:  "GET",
			Path:    "/docker/repositories/*name",
			Request: getDockerRepositories,
		},
	)

}

func injectServices() {
	// sampleService = new(service.SampleService)
	// httpsampleService = new(service.HTTPSampleService)

	dockerRegistryService = new(service.DockerRegistryService)
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
	r := dockerRegistryService.GetCatalog()

	c.JSON(http.StatusOK, r)
}

func getDockerRepositories(c *gin.Context) {
	repoName := c.Params.ByName("name")
	repoName = strings.Replace(repoName, "/", "", 1)

	if repoName == "" {
		r := dockerRegistryService.GetRepositories()
		c.JSON(http.StatusOK, r)
	} else {
		r := dockerRegistryService.GetRepository(repoName)
		if r.Tags == nil {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.JSON(http.StatusOK, r)
	}
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
