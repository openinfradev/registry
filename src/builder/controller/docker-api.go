package controller

import (
	"builder/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

var dockerService *service.DockerService

func init() {
	// inject service
	dockerService = new(service.DockerService)

	// docker build
	addRequestMapping(
		RequestMapper{
			Method:  "POST",
			Path:    "/docker/build",
			Request: buildDockerFile,
		},
	)

	// docker tag
	addRequestMapping(
		RequestMapper{
			Method:  "PATCH",
			Path:    "/docker/tag",
			Request: tagDockerImage,
		},
	)

	// docker push
	addRequestMapping(
		RequestMapper{
			Method:  "PUT",
			Path:    "/docker/push",
			Request: pushDockerImage,
		},
	)
}

// buildDockerFile
// @Summary docker build by dockerfile
// @Description docker build by dockerfile api
// @Name buildDockerFile
// @Produce  json
// @Router /docker/build [post]
// @Success 200 {object} model.BasicResult
func buildDockerFile(c *gin.Context) {
	// test arguments
	repoName := "exntu/sample2"
	dockerfilePath := "./sample"

	r := dockerService.Build(repoName, dockerfilePath)
	c.JSON(http.StatusOK, r)
}

// tagDockerImage
// @Summary docker image tag
// @Description docker image tag
// @Name tagDockerImage
// @Produce  json
// @Router /docker/tag [patch]
// @Success 200 {object} model.BasicResult
func tagDockerImage(c *gin.Context) {
	// test arguments
	repoName := "exntu/sample2"
	oldTag := "latest"
	newTag := "v100"

	r := dockerService.Tag(repoName, oldTag, newTag)
	c.JSON(http.StatusOK, r)
}

// pushDockerImage
// @Summary docker image push
// @Description docker image push
// @Name pushDockerImage
// @Produce  json
// @Router /docker/push [put]
// @Success 200 {object} model.BasicResult
func pushDockerImage(c *gin.Context) {
	// test arguments
	repoName := "exntu/sample2"
	tag := "v100"

	r := dockerService.Push(repoName, tag)
	c.JSON(http.StatusOK, r)
}
