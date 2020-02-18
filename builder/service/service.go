package service

import (
	"fmt"
	"github.com/openinfradev/registry-builder/builder/config"
	"github.com/openinfradev/registry-builder/builder/repository"
	"github.com/openinfradev/registry-builder/builder/util/logger"
	tokenutil "github.com/openinfradev/registry-builder/builder/util/token"
)

// InjectedServices is injection services
type InjectedServices struct {
	DockerService      *DockerService
	RegistryService    *RegistryService
	MinioService       *MinioService
	SecurityService    *SecurityService
	FileManager        *FileManager
	RegistryRepository *repository.RegistryRepository
}

var is *InjectedServices

func init() {
	is = &InjectedServices{
		DockerService:      new(DockerService),
		RegistryService:    new(RegistryService),
		MinioService:       new(MinioService),
		SecurityService:    new(SecurityService),
		FileManager:        new(FileManager),
		RegistryRepository: new(repository.RegistryRepository),
	}
}

// GetRegistryURL returns registry full url
func GetRegistryURL(path string) string {
	url := ""

	registryinfo := config.GetConfig().Registry

	if registryinfo.Insecure {
		url = "http://"
	} else {
		url = "https://"
	}
	url += registryinfo.Endpoint
	url += path

	// logger.DEBUG("service/service.go", "GetRegistryURL", url)

	return url
}

// Authorization is authorization
func Authorization(token string) bool {
	basicToken, err := tokenutil.DecodeBasicToken(token)
	if err != nil {
		logger.ERROR("service/service.go", "Authorization", err.Error())
	}
	logger.DEBUG("service/service.go", "Authorization", fmt.Sprintf("decoded token raw[%s]", basicToken.Raw))
	return true
}
