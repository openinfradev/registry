package controller

import (
	"builder/model"
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

	// registry manifest v1
	addRequestMapping(
		RequestMapper{
			Method:  "GET",
			Path:    "/registry/manifest-v1/*name",
			Request: getRegistryManifestV1,
		},
	)

	// registry manifest v2
	addRequestMapping(
		RequestMapper{
			Method:  "GET",
			Path:    "/registry/manifest-v2/*name",
			Request: getRegistryManifestV2,
		},
	)

}

// getRegistryCatalog
// @Summary docker registry catalog api
// @Description docker registry catalog api
// @Name getRegistryCatalog
// @Accept json
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
// @Accept json
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
// @Param tag query string false "Tag Name"
// @Accept json
// @Produce  json
// @Router /registry/repositories/{name} [delete]
// @Success 200 {object} model.BasicResult
func deleteRegistryRepository(c *gin.Context) {
	repoName := c.Params.ByName("name")
	repoName = strings.Replace(repoName, "/", "", 1)

	tag := c.Query("tag")
	r := &model.BasicResult{}
	if tag == "" {
		r = registryService.DeleteRepository(repoName)
	} else {
		r = registryService.DeleteRepositoryTag(repoName, tag)
	}
	c.JSON(http.StatusOK, r)
}

// getRegistryManifestV1
// @Summary docker registry manifest v1 api
// @Description docker registry manifest v1 api
// @Name getRegistryManifestV1
// @Param name path string true "Repository Name" default()
// @Param tag query string true "Tag Name"
// @Accept json
// @Produce  json
// @Router /registry/manifest-v1/{name} [get]
// @Success 200 {object} model.RegistryManifestV1
func getRegistryManifestV1(c *gin.Context) {
	repoName := c.Params.ByName("name")
	repoName = strings.Replace(repoName, "/", "", 1)

	tag := c.Query("tag")

	r := registryService.GetManifestV1(repoName, tag)
	c.JSON(http.StatusOK, r)
}

// getRegistryManifestV2
// @Summary docker registry manifest v2 api
// @Description docker registry manifest v2 api
// @Name getRegistryManifestV2
// @Param name path string true "Repository Name" default()
// @Param tag query string true "Tag Name"
// @Accept json
// @Produce  json
// @Router /registry/manifest-v2/{name} [get]
// @Success 200 {object} model.RegistryManifestV2
func getRegistryManifestV2(c *gin.Context) {
	repoName := c.Params.ByName("name")
	repoName = strings.Replace(repoName, "/", "", 1)

	tag := c.Query("tag")

	r := registryService.GetManifestV2(repoName, tag)
	c.JSON(http.StatusOK, r)
}
