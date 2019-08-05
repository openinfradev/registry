package service

import (
	"bufio"
	"builder/constant"
	"builder/util/logger"
	"fmt"
	"os/exec"
)

// DockerService is docker command relative service
type DockerService struct{}

// Build is docker building logs with dockerfile
func (d *DockerService) Build(repoName string, dockerfilePath string) *BasicResult {

	// needs using goroutine
	// and saving log line by line

	// async
	go buildJob(repoName, dockerfilePath)

	// only ok
	return &BasicResult{
		Code:    constant.ResultSuccess,
		Message: "",
	}
}

// Tag is image tagging
func (d *DockerService) Tag(repoName string, oldTag string, newTag string) *BasicResult {

	// needs using goroutine
	// and saving log line by line

	// sync
	ch := make(chan BasicResult, 1)
	go tagJob(ch, repoName, oldTag, newTag)
	r := <-ch

	return &r
}

// Push is docker image pushing
func (d *DockerService) Push(repoName string, tag string) *BasicResult {
	// needs using goroutine
	// and saving log line by line

	// async
	go pushJob(repoName, tag)

	// only ok
	return &BasicResult{
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

func tagJob(ch chan<- BasicResult, repoName string, oldTag string, newTag string) {
	logger.DEBUG("docker.go", fmt.Sprintf("tagJob [%s] [%s] to [%s]", repoName, oldTag, newTag))

	result := &BasicResult{}

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

	logger.DEBUG("docker.go", "buildJob end "+repoName)
}
