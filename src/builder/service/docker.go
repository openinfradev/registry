package service

import (
	"bufio"
	"builder/util/logger"
	"fmt"
	"os/exec"
)

// DockerService is docker command relative service
type DockerService struct{}

// Build is docker building logs with dockerfile
func (d *DockerService) Build(repoName string, dockerfilePath string) string {

	// needs using goroutine
	// and saving log line by line

	// async
	go buildJob(repoName, dockerfilePath)

	// only ok
	return `{"message":"ok"}`
}

// Tag is image tagging
func (d *DockerService) Tag(repoName string, oldTag string, newTag string) string {

	// needs using goroutine
	// and saving log line by line

	// sync
	ch := make(chan string, 1)
	go tagJob(ch, repoName, oldTag, newTag)
	r := <-ch

	return r
}

// Push is docker image pushing
func (d *DockerService) Push(repoName string, tag string) string {
	// needs using goroutine
	// and saving log line by line

	// async
	go pushJob(repoName, tag)

	// only ok
	return `{"message":"ok"}`
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

	logger.DEBUG("docker.go", fmt.Sprintf("pushJob end [%s:%s]", repoName, tag))
}

func tagJob(ch chan<- string, repoName string, oldTag string, newTag string) {
	logger.DEBUG("docker.go", fmt.Sprintf("tagJob [%s] [%s] to [%s]", repoName, oldTag, newTag))

	oldRepo := repoName + ":" + oldTag
	newRepo := basicinfo.RegistryEndpoint + "/" + repoName + ":" + newTag

	tag := exec.Command("docker", "tag", oldRepo, newRepo)

	err := tag.Run()
	if err != nil {
		logger.ERROR("docker.go", "tagJob is failed")
		ch <- `{"message":"error"}`
	} else {
		logger.DEBUG("docker.go", "tagJob is success")
		ch <- `{"message":"ok"}`
	}
}

func buildJob(repoName string, dockerfilePath string) {
	logger.DEBUG("docker.go", "buildJob start "+repoName)

	repoName = repoName + ":latest"
	build := exec.Command("docker", "build", "-t", repoName, dockerfilePath)

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
