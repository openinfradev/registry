package controller

import (
	"builder/model"
	"builder/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

var dockerService *service.DockerService

func init() {
	// inject service
	dockerService = new(service.DockerService)

	// docker build by file
	addRequestMapping(
		RequestMapper{
			Method:  "POST",
			Path:    "/docker/build/file",
			Request: buildByDockerFile,
		},
	)

	// docker build by git
	addRequestMapping(
		RequestMapper{
			Method:  "POST",
			Path:    "/docker/build/git",
			Request: buildByGitRepository,
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

// buildByDockerFile
// @Summary docker build by dockerfile
// @Description docker build by dockerfile api
// @Name buildByDockerFile
// @Accept json
// @Produce json
// @Router /docker/build/file [post]
// @Param contents body model.DockerBuildByFileParam true "Json Parameters (contents is base64 encoded)"
// @Success 200 {object} model.BasicResult
func buildByDockerFile(c *gin.Context) {

	var params *model.DockerBuildByFileParam
	c.BindJSON(&params)

	r := dockerService.BuildByDockerfile(params)

	c.JSON(http.StatusOK, r)
}

// buildByGitRepository
// @Summary docker build by git
// @Description docker build by git api
// @Name buildByGitRepository
// @Accept json
// @Produce json
// @Router /docker/build/git [post]
// @Param contents body model.DockerBuildByGitParam true "Json Parameters (userPW is base64 encoded)"
// @Success 200 {object} model.BasicResult
func buildByGitRepository(c *gin.Context) {

	var params *model.DockerBuildByGitParam
	c.BindJSON(&params)

	r := dockerService.BuildByGitRepository(params)
	c.JSON(http.StatusOK, r)
}

// tagDockerImage
// @Summary docker image tag
// @Description docker image tag
// @Name tagDockerImage
// @Produce  json
// @Router /docker/tag [patch]
// @Param contents body model.DockerTagParam true "Json Parameters"
// @Success 200 {object} model.BasicResult
func tagDockerImage(c *gin.Context) {

	var params *model.DockerTagParam
	c.BindJSON(&params)

	r := dockerService.Tag(params)
	c.JSON(http.StatusOK, r)
}

// pushDockerImage
// @Summary docker image push
// @Description docker image push
// @Name pushDockerImage
// @Produce  json
// @Router /docker/push [put]
// @Param contents body model.DockerPushParam true "Json Parameters"
// @Success 200 {object} model.BasicResult
func pushDockerImage(c *gin.Context) {

	var params *model.DockerPushParam
	c.BindJSON(&params)

	r := dockerService.Push(params)
	c.JSON(http.StatusOK, r)
}
