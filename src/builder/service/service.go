package service

import (
	"builder/config"
	"builder/repository"
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
	is = &InjectedServices {
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
