package controller

import (
	"builder/constant"
	"builder/model"
	"builder/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

var minioService *service.MinioService

func init() {

	// inject service
	minioService = new(service.MinioService)

	// create minio
	addRequestMapping(
		RequestMapper{
			Method:  "POST",
			Path:    "/minio",
			Request: createMinio,
		},
	)

	// delete minio
	addRequestMapping(
		RequestMapper{
			Method:  "DELETE",
			Path:    "/minio",
			Request: deleteMinio,
		},
	)

}

/*
	Request Mapping Functions
*/

// createMinio
// @Summary create minio
// @Description create minio
// @Name  createMinio
// @Produce  json
// @Router /minio [post]
// @Param contents body model.MinioParam true "Json Parameters"
// @Success 200 {object} model.MinioResult
func createMinio(c *gin.Context) {
	var params *model.MinioParam
	c.BindJSON(&params)

	r := minioService.CreateMinio(params)
	c.JSON(http.StatusOK, r)
}

// deleteMinio
// @Summary delete minio
// @Description delete minio
// @Name  deleteMinio
// @Produce  json
// @Router /minio [delete]
// @Param contents body model.MinioParam true "Json Parameters"
// @Success 200 {object} model.BasicResult
func deleteMinio(c *gin.Context) {
	var params *model.MinioParam
	c.BindJSON(&params)

	r := minioService.DeleteMinio(params.UserID)
	if r {
		c.JSON(http.StatusOK, &model.BasicResult{
			Code:    constant.ResultSuccess,
			Message: "",
		})
	} else {
		c.JSON(http.StatusOK, &model.BasicResult{
			Code:    constant.ResultFail,
			Message: "",
		})
	}
}
