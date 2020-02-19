package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func init() {

	// registry listener (hook)
	addRequestMapping(
		RequestMapper{
			Method:  "POST",
			Path:    "/listener",
			Request: listen,
		},
	)
}

// listener
// @Summary registry hook listen
// @Description registry hook listen
// @Name listener
// @Accept json
// @Produce json
// @Router /listener [post]
// @Success 200
func listen(c *gin.Context) {
	params := new(map[string]interface{})
	c.BindJSON(&params)

	is.WebhookService.Toss(params)

	c.AbortWithStatus(http.StatusOK)
}
