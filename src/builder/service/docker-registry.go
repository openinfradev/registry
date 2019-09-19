package service

import (
	"builder/constant"
	urlconst "builder/constant/url"
	"builder/model"
	"builder/util/logger"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// RegistryService is relative docker registry
type RegistryService struct{}

func init() {
}

// GetCatalog returns docker registry catalog
func (d *RegistryService) GetCatalog() *model.CatalogResult {

	catalogResult := &model.CatalogResult{}

	resp, err := http.Get(basicinfo.GetRegistryURL(urlconst.PathRegistryCatalog))
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
func (d *RegistryService) GetRepository(repoName string) *model.RepositoryResult {

	repositoryResult := &model.RepositoryResult{}

	path := fmt.Sprintf(urlconst.PathRegistryTagList, repoName)
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
func (d *RegistryService) GetRepositories() *model.RepositoriesResult {

	repositories := []model.RepositoryResult{}
	catalog := d.GetCatalog()
	for _, repoName := range catalog.Repositories {
		repository := d.GetRepository(repoName)
		// skipped tag which is null(or nil)
		if repository.Tags != nil {
			repositories = append(repositories, *repository)
		}
	}
	repositoriesResult := &model.RepositoriesResult{
		Repositories: repositories,
	}

	return repositoriesResult
}

// DeleteRepository is repository deleting
func (d *RegistryService) DeleteRepository(repoName string, tag string) *model.BasicResult {

	// get digest
	digest := d.GetDigest(repoName, tag)
	if digest == "" {
		return &model.BasicResult{
			Code:    constant.ResultFail,
			Message: "",
		}
	}

	// delete by digest
	path := fmt.Sprintf(urlconst.PathRegistryManifest, repoName, digest)
	req, err := http.NewRequest("DELETE", basicinfo.GetRegistryURL(path), nil)
	if err != nil {
		return &model.BasicResult{
			Code:    constant.ResultFail,
			Message: "",
		}
	}

	req.Header.Add("Accept", "application/vnd.docker.distribution.manifest.v2+json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return &model.BasicResult{
			Code:    constant.ResultFail,
			Message: "",
		}
	}
	defer resp.Body.Close()

	r, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return &model.BasicResult{
			Code:    constant.ResultFail,
			Message: "",
		}
	}

	// garbage collect (go-routine)
	// sync (???)
	ch := make(chan string, 1)
	go garbageCollectJob(ch)
	rr := <-ch
	logger.DEBUG("service/docker-registry.go", "DeleteRepository", rr)

	return &model.BasicResult{
		Code:    constant.ResultSuccess,
		Message: string(r),
	}
}

// GetDigest returns repository:tag digest
func (d *RegistryService) GetDigest(repoName string, tag string) string {

	path := fmt.Sprintf(urlconst.PathRegistryManifest, repoName, tag)
	req, err := http.NewRequest("GET", basicinfo.GetRegistryURL(path), nil)
	if err != nil {
		return ""
	}

	req.Header.Add("Accept", "application/vnd.docker.distribution.manifest.v2+json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	digest := resp.Header.Get("Docker-Content-Digest")
	return digest
}
