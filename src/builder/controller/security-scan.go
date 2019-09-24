package controller

import (
	"builder/service"
	"builder/util/logger"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

var securityService *service.SecurityService

func init() {
	// inject service
	securityService = new(service.SecurityService)

	// get layer scanning information
	addRequestMapping(
		RequestMapper{
			Method:  "GET",
			Path:    "/sescan/layer/:name",
			Request: getSecurityScanLayer,
		},
	)

	// patch layer scanning
	addRequestMapping(
		RequestMapper{
			Method:  "PATCH",
			Path:    "/sescan/repository/*name",
			Request: scanLayers,
		},
	)
}

// getSecurityScanLayer
// @Summary security scanning layer api
// @Description security scanning layer api
// @Name getSecurityScanLayer
// @Param name path string true "Layer(history) ID" default()
// @Accept json
// @Produce json
// @Router /sescan/layer/{name} [get]
// @Success 200 {object} model.SecurityScanLayer
func getSecurityScanLayer(c *gin.Context) {

	layerID := c.Params.ByName("name")

	r := securityService.GetLayer(layerID)

	c.JSON(http.StatusOK, r)
}

// scanLayers
// @Summary security scan api
// @Description security scan api
// @Name scanLayers
// @Param name path string true "Repository Name" default()
// @Param tag query string true "Tag Name"
// @Accept json
// @Produce json
// @Router /sescan/repository/{name} [patch]
// @Success 200 {object} model.BasicResult
func scanLayers(c *gin.Context) {
	repoName := c.Params.ByName("name")
	repoName = strings.Replace(repoName, "/", "", 1)

	tag := c.Query("tag")

	log := fmt.Sprintf("security scanning : repo[%s] tag[%s]", repoName, tag)
	logger.DEBUG("controller/security-scan.go", "scanLayers", log)

	r := securityService.Scan(repoName, tag)

	c.JSON(http.StatusOK, r)
}
