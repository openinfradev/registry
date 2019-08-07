package controller

import (
	"builder/service"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

var registryService *service.RegistryService

func init() {
	// inject service
	registryService = new(service.RegistryService)

	// registry catalog
	addRequestMapping(
		RequestMapper{
			Method:  "GET",
			Path:    "/registry/catalog",
			Request: getRegistryCatalog,
		},
	)

	// registry repositories tag list
	addRequestMapping(
		RequestMapper{
			Method:  "GET",
			Path:    "/registry/repositories/*name",
			Request: getRegistryRepositories,
		},
	)

	// registry repository delete
	addRequestMapping(
		RequestMapper{
			Method:  "DELETE",
			Path:    "/registry/repositories/*name",
			Request: deleteRegistryRepository,
		},
	)
}

// getRegistryCatalog
// @Summary docker registry catalog api
// @Description docker registry catalog api
// @Name getRegistryCatalog
// @Produce  json
// @Router /registry/catalog [get]
// @Success 200 {object} model.CatalogResult
func getRegistryCatalog(c *gin.Context) {
	r := registryService.GetCatalog()

	c.JSON(http.StatusOK, r)
}

// getRegistryRepositories
// @Summary docker registry repositories api
// @Description docker registry repositories api
// @Name getRegistryRepositories
// @Param name path string false "Repository Name" default()
// @Produce  json
// @Router /registry/repositories/{name} [get]
// @Success 200 {object} model.RepositoriesResult
// @Success 200 {object} model.RepositoryResult
func getRegistryRepositories(c *gin.Context) {
	repoName := c.Params.ByName("name")
	repoName = strings.Replace(repoName, "/", "", 1)

	if repoName == "" {
		r := registryService.GetRepositories()
		c.JSON(http.StatusOK, r)
	} else {
		r := registryService.GetRepository(repoName)
		if r.Tags == nil {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.JSON(http.StatusOK, r)
	}
}

// deleteRegistryRepository
// @Summary docker registry repository delete api
// @Description docker registry repository delete api
// @Name deleteRegistryRepository
// @Param name path string true "Repository Name" default()
// @Param tag query string true "Tag Name"
// @Produce  json
// @Router /registry/repositories/{name} [delete]
// @Success 200 {object} model.BasicResult
func deleteRegistryRepository(c *gin.Context) {
	repoName := c.Params.ByName("name")
	repoName = strings.Replace(repoName, "/", "", 1)

	tag := c.Query("tag")

	r := registryService.DeleteRepository(repoName, tag)
	c.JSON(http.StatusOK, r)
}
