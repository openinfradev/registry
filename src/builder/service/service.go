package service

import (
	"builder/util/logger"
)

// BasicInfo is defined service information
type BasicInfo struct {
	RegistryName     string
	RegistryInsecure bool
	RegistryEndpoint string
	TemporaryPath    string
}

var basicinfo *BasicInfo

// SetBasicInfo is setting service information
func SetBasicInfo(info *BasicInfo) {
	logger.DEBUG("service.go", "setting service information")

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

	logger.DEBUG("service.go GetRegistryURL", url)

	return url
}
