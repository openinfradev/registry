package service

import (
	"builder/constant"
	urlconst "builder/constant/url"
	"builder/model"
	"builder/util"
	"builder/util/logger"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// SecurityService is security scan service using clair
type SecurityService struct{}

func init() {
}

// GetLayer returns scanned layer vulnerabilities
func (s *SecurityService) GetLayer(layerID string) *model.SecurityScanLayer {

	layerScanResult := &model.SecurityScanLayer{}

	path := fmt.Sprintf(urlconst.SecurityScanLayer, basicinfo.ClairEndpoint, layerID)
	// layer features & vulnerabilities
	path += "?" + urlconst.SecurityScanLayerParam
	resp, err := http.Get(path)
	if err != nil {
		return layerScanResult
	}

	defer resp.Body.Close()

	statusCode := resp.StatusCode
	if statusCode == http.StatusNotFound {
		layerScanResult.Status = "queued"
		return layerScanResult
	}

	r, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return layerScanResult
	}

	// logger.DEBUG("service/security-scan.go", "GetLayer", string(r))

	// err ignore
	json.Unmarshal(r, &layerScanResult.Data)
	layerScanResult.Status = "scanned"

	return layerScanResult
}

// Scan is security scanning to clair. layer scan using manifests
func (s *SecurityService) Scan(repoName string, tag string) *model.BasicResult {

	ch := make(chan string, 1)
	go scanJob(ch, repoName, tag)

	return &model.BasicResult{
		Code:    constant.ResultSuccess,
		Message: "",
	}
}

func scanJob(ch chan<- string, repoName string, tag string) {

	registryService := new(RegistryService)
	manifest := registryService.GetManifestV1(repoName, tag)
	fsLayersMap := manifest["fsLayers"]
	historyMap := manifest["history"]

	fsLayers := []model.RegistryManifestV1Layer{}
	util.MapToStruct(fsLayersMap, &fsLayers)

	history := []model.RegistryManifestV1History{}
	util.MapToStruct(historyMap, &history)

	logger.DEBUG("service/security-scan.go", "scanJob", fmt.Sprintf("fsLayers[%d] history[%d]", len(fsLayers), len(history)))

	layerLen := len(fsLayers)
	// layers := []model.SecurityScanParam{}
	for i := 0; i < layerLen; i++ {
		// layer := &model.SecurityScanParam{}

	}
}
