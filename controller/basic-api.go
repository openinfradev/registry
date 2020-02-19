package controller

import (
	"bytes"
	"io/ioutil"

	"github.com/openinfradev/registry-builder/constant"
	"github.com/openinfradev/registry-builder/model"
	"github.com/openinfradev/registry-builder/util/logger"

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
			Method:  "POST",
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

	// dockerService := new(service.DockerService)
	// dockerService.Test("/home/linus/ngrinder", "/home/linus/ngrinder666")

	buf := make([]byte, 1024)
	num, _ := c.Request.Body.Read(buf)
	reqBody := string(buf[0:num])
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(reqBody)))

	buf2 := make([]byte, 1024)
	num2, _ := c.Request.Body.Read(buf2)
	reqBody2 := string(buf2[0:num2])

	logger.DEBUG("controller/registry-listener.go", "listen", reqBody2)

	c.JSON(http.StatusOK, buf2[0:num2])
}
