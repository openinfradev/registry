package service

import (
	"bufio"
	"builder/constant"
	"builder/model"
	"builder/util/logger"
	"encoding/base64"
	"fmt"
	"os/exec"
)

// DockerService is docker command relative service
type DockerService struct{}

var fileManager *FileManager

func init() {
	fileManager = new(FileManager)
}

// BuildByDockerfile is docker building by dockerfile
func (d *DockerService) BuildByDockerfile(repoName string, encodedContents string) *model.BasicResult {
	// needs using goroutine
	// and saving log line by line

	// decoding contents
	decoded, err := base64.StdEncoding.DecodeString(encodedContents)
	if err != nil {
		return &model.BasicResult{
			Code:    constant.ResultFail,
			Message: "contents isn't base64 encoded",
		}
	}

	path, err := fileManager.WriteDockerfile(string(decoded))
	if err != nil {
		return &model.BasicResult{
			Code:    constant.ResultFail,
			Message: "",
		}
	}

	return d.Build(repoName, path)
}

// BuildByGitRepository is docker building by git repository
func (d *DockerService) BuildByGitRepository(repoName string, gitRepo string, userID string, encodedUserPW string) *model.BasicResult {
	// needs using goroutine
	// and saving log line by line

	// decoding userPW
	decoded, err := base64.StdEncoding.DecodeString(encodedUserPW)
	if err != nil {
		return &model.BasicResult{
			Code:    constant.ResultFail,
			Message: "userPW isn't base64 encoded",
		}
	}

	// not using go-routine (not yet)
	// ch := make(chan string, 1)	// dirPath
	path, err := fileManager.PullGitRepository(gitRepo, userID, string(decoded))
	if err != nil {
		return &model.BasicResult{
			Code:    constant.ResultFail,
			Message: "",
		}
	}

	return d.Build(repoName, path)
}

// Build is docker building by file path
func (d *DockerService) Build(repoName string, dockerfilePath string) *model.BasicResult {

	// async
	go buildJob(repoName, dockerfilePath)

	// only ok
	return &model.BasicResult{
		Code:    constant.ResultSuccess,
		Message: "",
	}
}

// Tag is image tagging
func (d *DockerService) Tag(repoName string, oldTag string, newTag string) *model.BasicResult {

	// needs using goroutine
	// and saving log line by line

	// sync
	ch := make(chan model.BasicResult, 1)
	go tagJob(ch, repoName, oldTag, newTag)
	r := <-ch

	return &r
}

// Push is docker image pushing
func (d *DockerService) Push(repoName string, tag string) *model.BasicResult {
	// needs using goroutine
	// and saving log line by line

	// async
	go pushJob(repoName, tag)

	// only ok
	return &model.BasicResult{
		Code:    constant.ResultSuccess,
		Message: "",
	}
}

func pushJob(repoName string, tag string) {
	logger.DEBUG("docker.go", fmt.Sprintf("pushJob start [%s:%s]", repoName, tag))

	repoName = basicinfo.RegistryEndpoint + "/" + repoName + ":" + tag
	push := exec.Command("docker", "push", repoName)

	r := ""
	stdout, _ := push.StdoutPipe()
	push.Start()
	scanner := bufio.NewScanner(stdout)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		m := scanner.Text()
		r += m + "\n"
		logger.DEBUG("docker.go push ", m)
	}
	push.Wait()

	logger.DEBUG("docker.go", fmt.Sprintf("pushJob end [%s]", repoName))
}

func tagJob(ch chan<- model.BasicResult, repoName string, oldTag string, newTag string) {
	logger.DEBUG("docker.go", fmt.Sprintf("tagJob [%s] [%s] to [%s]", repoName, oldTag, newTag))

	result := &model.BasicResult{}

	oldRepo := repoName + ":" + oldTag
	newRepo := basicinfo.RegistryEndpoint + "/" + repoName + ":" + newTag

	tag := exec.Command("docker", "tag", oldRepo, newRepo)

	err := tag.Run()
	if err != nil {
		logger.ERROR("docker.go", "tagJob is failed")
		result.Code = constant.ResultFail
		result.Message = ""
		ch <- *result
	} else {
		logger.DEBUG("docker.go", "tagJob is success")
		result.Code = constant.ResultSuccess
		result.Message = ""
		ch <- *result
	}
}

func buildJob(repoName string, dockerfilePath string) {
	logger.DEBUG("docker.go", "buildJob start "+repoName)

	repoName = repoName + ":latest"
	build := exec.Command("docker", "build", "--no-cache", "-t", repoName, dockerfilePath)

	r := ""
	stdout, _ := build.StdoutPipe()
	build.Start()
	scanner := bufio.NewScanner(stdout)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		m := scanner.Text()
		r += m + "\n"
		logger.DEBUG("docker.go build", m)
	}
	build.Wait()

	// logger.DEBUG("docker.go", r)
	logger.DEBUG("docker.go", "buildJob end "+repoName)

	// path removeall
	fileManager.DeleteDirectory(dockerfilePath)
}

func garbageCollectJob(ch chan<- string) {
	logger.DEBUG("docker.go", "garbage collect start")

	gc := exec.Command("docker", "exec", basicinfo.RegistryName, "bin/registry", "garbage-collect", "/etc/docker/registry/config.yml")

	r := ""
	stdout, _ := gc.StdoutPipe()
	gc.Start()
	scanner := bufio.NewScanner(stdout)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		m := scanner.Text()
		r += m + "\n"
		logger.DEBUG("docker.go garbage collect", m)
	}
	gc.Wait()

	logger.DEBUG("docker.go", "garbage collect end")

	ch <- r
}
