package service

import (
	"builder/config"
	"builder/constant"
	"builder/constant/scope"
	urlconst "builder/constant/url"
	"builder/model"
	"builder/util/logger"
	tokenutil "builder/util/token"
	"encoding/json"
	"errors"
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

	token, err := d.Authorization(&scope.Scope{
		Type:     scope.TypeRegistry,
		Resource: scope.ResourceCatalog,
		Action:   scope.ActionWildCard,
	})
	if err != nil {
		return catalogResult
	}
	req, err := http.NewRequest("GET", GetRegistryURL(urlconst.PathRegistryCatalog), nil)
	if err != nil {
		return catalogResult
	}

	req.Header.Add("Authorization", token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.ERROR("service/docker-registry.go", "GetCatalog", err.Error())
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

	token, err := d.Authorization(&scope.Scope{
		Type:     scope.TypeRepository,
		Resource: repoName,
		Action:   scope.ActionPull,
	})
	if err != nil {
		return repositoryResult
	}
	req, err := http.NewRequest("GET", GetRegistryURL(path), nil)
	if err != nil {
		return repositoryResult
	}

	req.Header.Add("Authorization", token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return repositoryResult
	}

	defer resp.Body.Close()

	r, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return repositoryResult
	}

	// logger.DEBUG("service/docker-registry.go", "GetRepository", string(r))

	// err ignore
	json.Unmarshal(r, repositoryRaw)
	repositoryResult.Name = repositoryRaw.Name
	tags := []model.TagResult{}
	if repositoryRaw.Tags != nil && len(repositoryRaw.Tags) > 0 {
		for _, tagName := range repositoryRaw.Tags {
			digest := d.GetDigest(repoName, tagName)
			// logger.DEBUG("service/docker-registry.go", "GetRepository", fmt.Sprintf("%s:%s : digest [%s]", repoName, tagName, digest))
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

	logger.DEBUG("service/docker-registry.go", "DeleteRepositoryTag", fmt.Sprintf("repoName[%s] tag[%s]", repoName, tag))

	token, err := d.Authorization(&scope.Scope{
		Type:     scope.TypeRepository,
		Resource: repoName,
		Action:   scope.ActionWildCard,
	})
	if err != nil {
		return &model.BasicResult{
			Code:    constant.ResultFail,
			Message: "",
		}
	}

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
	req, err := http.NewRequest("DELETE", GetRegistryURL(path), nil)
	if err != nil {
		return &model.BasicResult{
			Code:    constant.ResultFail,
			Message: "",
		}
	}

	req.Header.Add("Accept", "application/vnd.docker.distribution.manifest.v2+json")
	req.Header.Add("Authorization", token)

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

	// logger.DEBUG("service/docker-registry.go", "DeleteRepositoryTag", string(r))

	// garbage collect (go-routine)
	// deprecated : has a problem
	// ch := make(chan string, 1)
	// go garbageCollectJob(ch)
	// rr := <-ch
	// logger.DEBUG("service/docker-registry.go", "DeleteRepository", rr)

	return &model.BasicResult{
		Code:    constant.ResultSuccess,
		Message: string(r),
	}
}

// GetDigest returns repository:tag digest
func (d *RegistryService) GetDigest(repoName string, tag string) string {

	token, err := d.Authorization(&scope.Scope{
		Type:     scope.TypeRepository,
		Resource: repoName,
		Action:   scope.ActionPull,
	})
	if err != nil {
		return ""
	}

	path := fmt.Sprintf(urlconst.PathRegistryManifest, repoName, tag)
	req, err := http.NewRequest("GET", GetRegistryURL(path), nil)
	if err != nil {
		return ""
	}

	req.Header.Add("Accept", "application/vnd.docker.distribution.manifest.v2+json")
	req.Header.Add("Authorization", token)

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

	token, err := d.Authorization(&scope.Scope{
		Type:     scope.TypeRepository,
		Resource: repoName,
		Action:   scope.ActionPull,
	})
	if err != nil {
		return nil
	}

	path := fmt.Sprintf(urlconst.PathRegistryManifest, repoName, tag)
	req, err := http.NewRequest("GET", GetRegistryURL(path), nil)
	if err != nil {
		return nil
	}

	req.Header.Add("Accept", "application/vnd.docker.distribution.manifest.v1+prettyjws")
	req.Header.Add("Authorization", token)

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

	token, err := d.Authorization(&scope.Scope{
		Type:     scope.TypeRepository,
		Resource: repoName,
		Action:   scope.ActionPull,
	})
	if err != nil {
		return nil
	}

	path := fmt.Sprintf(urlconst.PathRegistryManifest, repoName, tag)
	req, err := http.NewRequest("GET", GetRegistryURL(path), nil)
	if err != nil {
		return nil
	}

	req.Header.Add("Accept", "application/vnd.docker.distribution.manifest.v2+json")
	req.Header.Add("Authorization", token)

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

// Authorization returns docker registry authorization token
func (d *RegistryService) Authorization(scope *scope.Scope) (string, error) {
	registryinfo := config.GetConfig().Registry
	if registryinfo.Auth == "" {
		return "", errors.New("Authorization endpoint argument is empty")
	}

	path := fmt.Sprintf(registryinfo.Auth+"?scope=%s&service=%s", scope.String(), registryinfo.Name)
	// logger.DEBUG("service/docker-registry.go", "Authorization", path)
	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return "", err
	}

	req.Header.Add("Authorization", tokenutil.BuilderBasicToken())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.ERROR("service/docker-registry.go", "Authorization", err.Error())
		return "", err
	}
	defer resp.Body.Close()

	r, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	m := make(map[string]string)
	json.Unmarshal(r, &m)

	// logger.DEBUG("service/docker-registry.go", "Authorization", m["token"])

	return fmt.Sprintf("Bearer %s", m["token"]), nil
}
