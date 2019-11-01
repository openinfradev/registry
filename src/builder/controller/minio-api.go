package controller

import (
	"builder/model"

	"github.com/gin-gonic/gin"
)

func init() {

	// create minio
	addRequestMapping(
		RequestMapper{
			Method:  "POST",
			Path:    "/minio",
			Request: createMinio,
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
// @Success 200 {object} model.BasicResult
func createMinio(c *gin.Context) {
	var params *model.MinioParam
	c.BindJSON(&params)

	//
}
