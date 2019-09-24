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

	repositoryRaw := &model.RepositoryRaw{}
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
	json.Unmarshal(r, repositoryRaw)
	repositoryResult.Name = repositoryRaw.Name
	tags := []model.TagResult{}
	if repositoryRaw.Tags != nil && len(repositoryRaw.Tags) > 0 {
		for _, tagName := range repositoryRaw.Tags {
			digest := d.GetDigest(repoName, tagName)
			logger.DEBUG("service/docker-registry.go", "GetRepository", fmt.Sprintf("%s:%s : digest [%s]", repoName, tagName, digest))
			tag := &model.TagResult{
				Name:   tagName,
				Digest: digest,
			}
			tags = append(tags, *tag)
		}
	}
	repositoryResult.Tags = tags

	return repositoryResult
}

// GetRepositories returns repositories included tags
func (d *RegistryService) GetRepositories() *model.RepositoriesResult {

	repositories := []model.RepositoryResult{}
	catalog := d.GetCatalog()
	for _, repoName := range catalog.Repositories {
		repository := d.GetRepository(repoName)
		// skipped tag which is null(or nil)
		if repository.Tags != nil && len(repository.Tags) > 0 {
			repositories = append(repositories, *repository)
		}
	}
	repositoriesResult := &model.RepositoriesResult{
		Repositories: repositories,
	}

	return repositoriesResult
}

// DeleteRepository is repository deleting
func (d *RegistryService) DeleteRepository(repoName string) *model.BasicResult {
	repo := d.GetRepository(repoName)
	if repo.Tags == nil || len(repo.Tags) < 1 {
		return &model.BasicResult{
			Code:    constant.ResultFail,
			Message: "",
		}
	}
	for _, tag := range repo.Tags {
		d.DeleteRepositoryTag(repoName, tag.Name)
	}
	return &model.BasicResult{
		Code:    constant.ResultSuccess,
		Message: "",
	}
}

// DeleteRepositoryTag is repository tag deleting
func (d *RegistryService) DeleteRepositoryTag(repoName string, tag string) *model.BasicResult {

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

// GetManifestV1 returns registry manifest v1
func (d *RegistryService) GetManifestV1(repoName string, tag string) map[string]interface{} {

	path := fmt.Sprintf(urlconst.PathRegistryManifest, repoName, tag)
	req, err := http.NewRequest("GET", basicinfo.GetRegistryURL(path), nil)
	if err != nil {
		return nil
	}

	req.Header.Add("Accept", "application/vnd.docker.distribution.manifest.v1+prettyjws")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	r, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil
	}

	// err ignore
	result := make(map[string]interface{})
	json.Unmarshal(r, &result)

	return result
}

// GetManifestV2 returns registry manifest v2
func (d *RegistryService) GetManifestV2(repoName string, tag string) map[string]interface{} {

	path := fmt.Sprintf(urlconst.PathRegistryManifest, repoName, tag)
	req, err := http.NewRequest("GET", basicinfo.GetRegistryURL(path), nil)
	if err != nil {
		return nil
	}

	req.Header.Add("Accept", "application/vnd.docker.distribution.manifest.v2+json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	r, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil
	}

	// err ignore
	result := make(map[string]interface{})
	json.Unmarshal(r, &result)

	return result
}
