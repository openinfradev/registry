package service

import (
	"builder/constant"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// RegistryService is relative docker registry
type RegistryService struct{}

// CatalogResult is registry catalog result
type CatalogResult struct {
	Repositories []string `json:"repositories"`
}

// RepositoryResult is registry repository result
type RepositoryResult struct {
	Name string   `json:"name"`
	Tags []string `json:"tags"`
}

// RepositoriesResult is regitry repositories result
type RepositoriesResult struct {
	Repositories []RepositoryResult `json:"repositories"`
}

// GetCatalog returns docker registry catalog
func (d *RegistryService) GetCatalog() *CatalogResult {
	// needs admin logon
	// needs token

	resp, err := http.Get(basicinfo.GetRegistryURL(constant.PathRegistryCatalog))
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	catalogResult := &CatalogResult{}

	r, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return catalogResult
	}

	// err ignore
	json.Unmarshal(r, catalogResult)

	return catalogResult
}

// GetRepository returns repository included tags
func (d *RegistryService) GetRepository(repoName string) *RepositoryResult {
	// needs admin logon
	// needs token

	path := fmt.Sprintf(constant.PathRegistryTagList, repoName)
	resp, err := http.Get(basicinfo.GetRegistryURL(path))
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	repositoryResult := &RepositoryResult{}

	r, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return repositoryResult
	}

	// err ignore
	json.Unmarshal(r, repositoryResult)

	return repositoryResult
}

// GetRepositories returns repositories included tags
func (d *RegistryService) GetRepositories() *RepositoriesResult {
	// needs admin logon
	// needs token

	repositories := []RepositoryResult{}
	catalog := d.GetCatalog()
	for _, repoName := range catalog.Repositories {
		repository := d.GetRepository(repoName)
		repositories = append(repositories, *repository)
	}
	repositoriesResult := &RepositoriesResult{
		Repositories: repositories,
	}

	return repositoriesResult
}
