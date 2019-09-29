package controller

import (
	"builder/constant"
	"builder/model"
	"builder/util"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func init() {

	// health
	addRequestMapping(
		RequestMapper{
			Method:  "GET",
			Path:    "/health",
			Request: health,
		},
	)

	// test
	addRequestMapping(
		RequestMapper{
			Method:  "GET",
			Path:    "/test",
			Request: test,
		},
	)
}

/*
	Request Mapping Functions
*/

// health
// @Summary health check api
// @Description builder의 health를 체크할 목적의 api
// @Name health
// @Produce  json
// @Router /health [get]
// @Success 200 {object} model.BasicResult
func health(c *gin.Context) {
	c.JSON(http.StatusOK, &model.BasicResult{
		Code:    constant.ResultSuccess,
		Message: "taco-registry-builder is healthy",
	})
}

// test
// @Summary test api
// @Description test
// @Name test
// @Produce  json
// @Router /test [get]
// @Success 200
func test(c *gin.Context) {

	sample := "linus@github.com/repo/aaaa.git"
	gitRepo := util.ExtractGitRepositoryURL(sample)
	c.JSON(http.StatusOK, &model.BasicResult{
		Code:    constant.ResultSuccess,
		Message: fmt.Sprintf("protocol[%s] url[%s]", gitRepo.Protocol, gitRepo.URL),
	})
}
