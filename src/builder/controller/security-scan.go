package controller

import (
	"builder/service"
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

	// get layer scanning information by repo name & tag
	addRequestMapping(
		RequestMapper{
			Method:  "GET",
			Path:    "/sescan/repository/*name",
			Request: getSecurityScanLayerByRepo,
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

// getSecurityScanLayerByRepo
// @Summary security scanning layer api
// @Description security scanning layer api
// @Name getSecurityScanLayerByRepo
// @Param name path string true "Repository Name" default()
// @Param tag query string true "Tag Name"
// @Accept json
// @Produce  json
// @Router /sescan/repository/{name} [get]
// @Success 200 {object} model.SecurityScanLayer
func getSecurityScanLayerByRepo(c *gin.Context) {

	repoName := c.Params.ByName("name")
	repoName = strings.Replace(repoName, "/", "", 1)

	tag := c.Query("tag")

	r := securityService.GetLayerByRepo(repoName, tag)

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

	r := securityService.Scan(repoName, tag)

	c.JSON(http.StatusOK, r)
}
