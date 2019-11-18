package service

import (
	"builder/util/logger"
	"builder/repository"
)

// BasicInfo is defined service information
type BasicInfo struct {
	RegistryName     string
	RegistryInsecure bool
	RegistryEndpoint string
	TemporaryPath    string
	RedisEndpoint    string
	ClairEndpoint    string
	AuthURL          string
	ServiceDomain    string
	ServicePort      string
	MinioData        string
	MinioDomain      string
}

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

var basicinfo *BasicInfo

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

// SetBasicInfo is setting service information
func SetBasicInfo(info *BasicInfo) {
	logger.INFO("service/service.go", "SetBasicInfo", "setting service information")

	basicinfo = info
}

// GetRegistryURL returns registry full url
func (b *BasicInfo) GetRegistryURL(path string) string {
	url := ""
	if basicinfo.RegistryInsecure {
		url = "http://"
	} else {
		url = "https://"
	}
	url += basicinfo.RegistryEndpoint
	url += path

	// logger.DEBUG("service/service.go", "GetRegistryURL", url)

	return url
}
