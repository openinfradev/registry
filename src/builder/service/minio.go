package service

import (
	"builder/constant"
	"builder/constant/minio"
	"builder/model"
	"builder/util/logger"
	"encoding/base64"
	"fmt"
)

// MinioService is minio service
type MinioService struct{}

// CreateMinio is creating minio container by user
func (m *MinioService) CreateMinio(params *model.MinioParam) *model.BasicResult {

	// 0. decoded password
	decoded, err := base64.StdEncoding.DecodeString(params.UserPW)
	if err != nil {
		return &model.BasicResult{
			Code:    constant.ResultFail,
			Message: "userPW isn't base64 encoded",
		}
	}

	// 1. make directory (if not exists)
	fileManager := new(FileManager)
	mntPath := fileManager.MakeDirectory(basicinfo.MinioDirectory, params.UserID)

	// 2. pulling minio docker image
	pullParam := &model.DockerPullParam{
		Name: minio.MinioImageName,
		Tag:  minio.MinioImageTag,
	}
	dockerService := new(DockerService)
	r := dockerService.Pull(pullParam, false, true)
	if r.Code != constant.ResultSuccess {
		return r
	}

	// 3. minio port process

	// 4. run minio container
	/*
		docker run \
		-d --restart=always \
		-p 9000:9000 --name minio \
		-v /home/ubuntu/minio/data:/data \
		-e "MINIO_ACCESS_KEY=exntu" \
		-e "MINIO_SECRET_KEY=exntu123!" \
		minio/minio:latest server /data

		"docker run -d --restart=always -p %d:%d --name taco-minio-%s -v %s:/data -e \"MINIO_ACCESS_KEY=%s\" -e \"MINIO_SECRET_KEY=%s\" %s:%s server /data"
	*/
	logger.DEBUG("service/minio.go", "CreateMinio", fmt.Sprintf("%s %s %s", params.UserID, decoded, mntPath))

	return &model.BasicResult{
		Code:    constant.ResultSuccess,
		Message: "",
	}
}
