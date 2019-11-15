package service

import (
	"builder/constant"
	"builder/constant/minio"
	"builder/model"
	"builder/repository"
	"builder/util/logger"
	"encoding/base64"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

// MinioService is minio service
type MinioService struct{}

// ExistsContainer returns container exists
func (m *MinioService) ExistsContainer(userID string) bool {
	stdout, err := exec.Command("/bin/sh", "-c", fmt.Sprintf(minio.MinioContainerExistsTemplate, userID)).Output()
	if err != nil || strings.TrimSpace(string(stdout)) == "" {
		return false
	}
	return true
}

// CreateMinio is creating minio container by user
func (m *MinioService) CreateMinio(params *model.MinioParam) *model.MinioResult {

	// 0. decoded password
	decoded, err := base64.StdEncoding.DecodeString(params.UserPW)
	if err != nil {
		return &model.MinioResult{
			BasicResult: model.BasicResult{
				Code:    constant.ResultFail,
				Message: "userPW isn't base64 encoded",
			},
		}
	}

	// 1. make directory (if not exists)
	fileManager := new(FileManager)
	fileManager.MakeDirectory(minio.MinioDataPath, params.UserID)
	mountPath := fmt.Sprintf("%s/%s", basicinfo.MinioData, params.UserID)

	// 2. pulling minio docker image
	pullParam := &model.DockerPullParam{
		Name: minio.MinioImageName,
		Tag:  minio.MinioImageTag,
	}
	dockerService := new(DockerService)
	r := dockerService.Pull(pullParam, false, true)
	if r.Code != constant.ResultSuccess {
		return &model.MinioResult{
			BasicResult: *r,
		}
	}

	// 3-0. clean
	exists := m.ExistsContainer(params.UserID)
	if exists {
		m.DeleteMinio(params.UserID)
	}

	// 3. minio port process
	registryRepository := new(repository.RegistryRepository)
	registryRepository.CreatePortTableIfExists()
	topPort := registryRepository.GetTopPort()
	topPort++

	// 4. run minio container
	run := exec.Command("/bin/sh", "-c", fmt.Sprintf(minio.MinioDockerRunTemplate, topPort, params.UserID, mountPath, params.UserID, decoded, minio.MinioImageName, minio.MinioImageTag))
	err = run.Run()
	if err != nil {
		logger.ERROR("service/minio.go", "CreateMinio", err.Error())
		return &model.MinioResult{
			BasicResult: model.BasicResult{
				Code:    constant.ResultFail,
				Message: "Failed to create minio container",
			},
		}
	}

	// 5. check container alive
	exists = m.ExistsContainer(params.UserID)
	if !exists {
		logger.ERROR("service/minio.go", "CreateMinio", "Failed to run minio container")
		return &model.MinioResult{
			BasicResult: model.BasicResult{
				Code:    constant.ResultFail,
				Message: "Failed to run minio container",
			},
		}
	}

	// 6. insert new port
	registryRepository.InsertPort(topPort)

	// logger.DEBUG("service/minio.go", "CreateMinio", fmt.Sprintf("%s %s %s", params.UserID, decoded, mountPath))

	return &model.MinioResult{
		BasicResult: model.BasicResult{
			Code:    constant.ResultSuccess,
			Message: "",
		},
		Domain: basicinfo.MinioDomain,
		Port:   topPort,
	}
}

// DeleteMinio is deleting minio container
func (m *MinioService) DeleteMinio(userID string) bool {
	stdout, _ := exec.Command("/bin/sh", "-c", fmt.Sprintf(minio.MinioGetContainerPortTemplate, userID)).Output()
	port := strings.TrimSpace(string(stdout))
	if port != "" {
		// delete temporary port
		logger.DEBUG("service/minio.go", "DeleteMinio", fmt.Sprintf("Deletion target container port [%s]", port))
		registryRepository := new(repository.RegistryRepository)
		iport, _ := strconv.Atoi(port)
		registryRepository.DeletePort(iport)
	}

	// remove minio container
	_, err := exec.Command("/bin/sh", "-c", fmt.Sprintf(minio.MinioRemoveContainerTemplate, userID)).Output()
	if err != nil {
		logger.ERROR("service/minio.go", "DeleteMinio", err.Error())
		return false
	}

	// delete user's directory
	fileManager := new(FileManager)
	fileManager.DeleteDirectory(fmt.Sprintf("%s/%s", minio.MinioDataPath, userID))

	return true
}
