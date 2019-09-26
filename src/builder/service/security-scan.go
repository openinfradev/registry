package service

import (
	"builder/constant"
	"builder/constant/scope"
	urlconst "builder/constant/url"
	"builder/model"
	"builder/util"
	"builder/util/logger"
	"bytes"
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

// GetLayerByRepo returns scanned layer vulnerabilities
func (s *SecurityService) GetLayerByRepo(repoName string, tag string) *model.SecurityScanLayer {

	registryService := new(RegistryService)
	manifest := registryService.GetManifestV1(repoName, tag)
	if manifest == nil {
		logger.ERROR("service/security-scan.go", "GetLayerByRepo", fmt.Sprintf("Not exists manifest [%s:%s]", repoName, tag))
		return nil
	}
	historyMap := manifest["history"]
	history := []model.RegistryManifestV1History{}
	util.MapToStruct(historyMap, &history)

	if len(history) > 0 {
		h0 := &model.RegistryManifestV1HistoryValue{}
		json.Unmarshal([]byte(history[0].V1Compatibility), h0)
		return s.GetLayer(h0.ID)
	}
	return nil
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
	logger.DEBUG("service/security-scan.go", "scanJob", fmt.Sprintf("start [%s:%s]", repoName, tag))

	// 0. authorization token
	registryService := new(RegistryService)
	token, err := registryService.Authorization(&scope.Scope{
		Type:     scope.TypeRepository,
		Resource: repoName,
		Action:   scope.ActionPull,
	})
	if err != nil {
		ch <- constant.ResultFail
		return
	}

	// 1. get manifest
	manifest := registryService.GetManifestV1(repoName, tag)
	if manifest == nil {
		logger.ERROR("service/security-scan.go", "scanJob", fmt.Sprintf("Not exists manifest [%s:%s]", repoName, tag))
		ch <- constant.ResultFail
		return
	}
	fsLayersMap := manifest["fsLayers"]
	historyMap := manifest["history"]

	// 2. parse fsLayers, history
	fsLayers := []model.RegistryManifestV1Layer{}
	util.MapToStruct(fsLayersMap, &fsLayers)

	history := []model.RegistryManifestV1History{}
	util.MapToStruct(historyMap, &history)

	// 3. make scan layer parameters
	layerLen := len(fsLayers)
	params := []model.SecurityScanParam{}
	for i := 0; i < layerLen; i++ {
		cs := history[i].V1Compatibility
		c := &model.RegistryManifestV1HistoryValue{}
		json.Unmarshal([]byte(cs), c)

		path := fmt.Sprintf(basicinfo.GetRegistryURL(urlconst.PathRegistryBlobs), repoName, fsLayers[i].BlobSum)

		layer := &model.SecurityScanLayerParam{}
		layer.Name = c.ID
		layer.ParentName = c.Parent
		layer.Path = path
		layer.Format = "Docker"

		headers := &model.SecurityScanLayerHeaderParam{
			Authorization: token,
		}
		layer.Headers = headers

		param := &model.SecurityScanParam{
			Layer: layer,
		}

		params = append(params, *param)
	}

	// 4. scan layer loop
	if len(params) > 0 {
		params = hierarchySort(params)
		for _, param := range params {
			requestScan(&param)
		}
		logger.DEBUG("service/security-scan.go", "scanJob", fmt.Sprintf("end [%s:%s] %d layers [%s]", repoName, tag, len(params), params[len(params)-1].Layer.Name))
		ch <- constant.ResultSuccess
	} else {
		logger.DEBUG("service/security-scan.go", "scanJob", fmt.Sprintf("end [%s:%s] not exists layers", repoName, tag))
		ch <- constant.ResultFail
	}
}

func requestScan(param *model.SecurityScanParam) {
	b, _ := json.Marshal(param)
	buff := bytes.NewBuffer(b)
	path := fmt.Sprintf(urlconst.SecurityScan, basicinfo.ClairEndpoint)
	resp, err := http.Post(path, "application/json", buff)
	if err != nil {
		logger.ERROR("service/security-scan.go", "requestScan", err.Error())
		return
	}

	defer resp.Body.Close()

	// It's not necessary.
	r, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.ERROR("service/security-scan.go", "requestScan", err.Error())
		return
	}
	logger.DEBUG("service/security-scan.go", "requestScan", string(r))
}

func hierarchySort(raw []model.SecurityScanParam) []model.SecurityScanParam {

	var root model.SecurityScanParam
	for _, p := range raw {
		if p.Layer.ParentName == "" {
			root = p
			break
		}
	}

	dist := []model.SecurityScanParam{}
	dist = append(dist, root)

	target := root
	loop := true
	for loop {
		exist := false
		for _, p := range raw {
			if target.Layer.Name == p.Layer.ParentName {
				target = p
				dist = append(dist, p)
				exist = true
				break
			}
		}
		loop = exist
	}
	return dist
}
