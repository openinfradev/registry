package service

import (
	"builder/constant"
	"builder/util/logger"
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

	catalogResult := &CatalogResult{}

	resp, err := http.Get(basicinfo.GetRegistryURL(constant.PathRegistryCatalog))
	if err != nil {
		return catalogResult
	}

	defer resp.Body.Close()

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

	repositoryResult := &RepositoryResult{}

	path := fmt.Sprintf(constant.PathRegistryTagList, repoName)
	resp, err := http.Get(basicinfo.GetRegistryURL(path))
	if err != nil {
		return repositoryResult
	}

	defer resp.Body.Close()

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
		// skipped tag which is null(or nil)
		if repository.Tags != nil {
			repositories = append(repositories, *repository)
		}
	}
	repositoriesResult := &RepositoriesResult{
		Repositories: repositories,
	}

	return repositoriesResult
}

// DeleteRepository is repository deleting
func (d *RegistryService) DeleteRepository(repoName string, tag string) *BasicResult {

	// get digest
	path := fmt.Sprintf(constant.PathRegistryManifest, repoName, tag)
	req, err := http.NewRequest("GET", basicinfo.GetRegistryURL(path), nil)
	if err != nil {
		return &BasicResult{
			Code:    constant.ResultFail,
			Message: "",
		}
	}

	req.Header.Add("Accept", "application/vnd.docker.distribution.manifest.v2+json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return &BasicResult{
			Code:    constant.ResultFail,
			Message: "",
		}
	}
	defer resp.Body.Close()

	digest := resp.Header.Get("Docker-Content-Digest")
	if digest == "" {
		return &BasicResult{
			Code:    constant.ResultFail,
			Message: "",
		}
	}

	// delete by digest
	path = fmt.Sprintf(constant.PathRegistryManifest, repoName, digest)
	req, err = http.NewRequest("DELETE", basicinfo.GetRegistryURL(path), nil)
	if err != nil {
		return &BasicResult{
			Code:    constant.ResultFail,
			Message: "",
		}
	}

	req.Header.Add("Accept", "application/vnd.docker.distribution.manifest.v2+json")

	client = &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		return &BasicResult{
			Code:    constant.ResultFail,
			Message: "",
		}
	}
	defer resp.Body.Close()

	r, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return &BasicResult{
			Code:    constant.ResultFail,
			Message: "",
		}
	}

	// garbage collect (go-routine)
	// sync (???)
	ch := make(chan string, 1)
	go garbageCollectJob(ch)
	rr := <-ch
	logger.DEBUG("docker-registry.go", rr)

	return &BasicResult{
		Code:    constant.ResultSuccess,
		Message: string(r),
	}
}
