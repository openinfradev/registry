package service

import (
	"builder/config"
	"builder/constant"
	"builder/constant/minio"
	"builder/model"
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

	minioinfo := config.GetConfig().Minio

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
	is.FileManager.MakeDirectory(minio.MinioDataPath, params.UserID)
	mountPath := fmt.Sprintf("%s/%s", minioinfo.Data, params.UserID)

	// 2. pulling minio docker image
	pullParam := &model.DockerPullParam{
		Name: minio.MinioImageName,
		Tag:  minio.MinioImageTag,
	}
	r := is.DockerService.Pull(pullParam, false, true)
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
	availablePort, err := is.RegistryRepository.GetAvailablePort()
	if err != nil {
		return &model.MinioResult{
			BasicResult: model.BasicResult{
				Code:    constant.ResultFail,
				Message: err.Error(),
			},
		}
	}

	// 4. run minio container
	run := exec.Command("/bin/sh", "-c", fmt.Sprintf(minio.MinioDockerRunTemplate, availablePort, params.UserID, mountPath, params.UserID, decoded, minio.MinioImageName, minio.MinioImageTag))
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

	// 6. delete port
	is.RegistryRepository.DeletePort(availablePort)

	// logger.DEBUG("service/minio.go", "CreateMinio", fmt.Sprintf("%s %s %s", params.UserID, decoded, mountPath))

	return &model.MinioResult{
		BasicResult: model.BasicResult{
			Code:    constant.ResultSuccess,
			Message: "",
		},
		Domain: minioinfo.Domain,
		Port:   availablePort,
	}
}

// DeleteMinio is deleting minio container
func (m *MinioService) DeleteMinio(userID string) bool {
	stdout, _ := exec.Command("/bin/sh", "-c", fmt.Sprintf(minio.MinioGetContainerPortTemplate, userID)).Output()
	port := strings.TrimSpace(string(stdout))
	if port != "" {
		// delete temporary port
		logger.DEBUG("service/minio.go", "DeleteMinio", fmt.Sprintf("Deletion target container port [%s]", port))
		iport, _ := strconv.Atoi(port)
		is.RegistryRepository.InsertPort(iport)
	}

	// remove minio container
	_, err := exec.Command("/bin/sh", "-c", fmt.Sprintf(minio.MinioRemoveContainerTemplate, userID)).Output()
	if err != nil {
		logger.ERROR("service/minio.go", "DeleteMinio", err.Error())
		return false
	}

	// delete user's directory
	is.FileManager.DeleteDirectory(fmt.Sprintf("%s/%s", minio.MinioDataPath, userID))

	return true
}
