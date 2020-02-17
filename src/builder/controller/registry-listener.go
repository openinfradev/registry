package controller

import (
	"bytes"
	"io/ioutil"
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
	buf := make([]byte, 1024)
	num, _ := c.Request.Body.Read(buf)
	reqBody := string(buf[0:num])
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(reqBody)))

	buf2 := make([]byte, 1024)
	num2, _ := c.Request.Body.Read(buf2)

	is.WebhookService.Toss(buf2[0:num2])

	c.AbortWithStatus(http.StatusOK)
}
