package controller

import (
	"builder/constant"
	"builder/model"
	"builder/util/logger"
	tokenutil "builder/util/token"
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
	t := "Basic ZXhudHU6ZXhudHUxMjM="
	bt, _ := tokenutil.ParseBasicToken(t)
	logger.DEBUG("controller/basic-api.go", "health", fmt.Sprintf("raw[%s] username[%s] password[%s]", bt.Raw, bt.Username, bt.Password))

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
	// t := taco.ParseLog("build-id", 1, "Step    91111114444440   : EXPOSE 22")

	// c.JSON(http.StatusOK, t)
}
