package service

import (
	"bufio"
	"builder/util/logger"
	"os/exec"
)

// DockerService is docker command relative service
type DockerService struct{}

// Build returns docker building logs with dockerfile
func (d *DockerService) Build(repoName string, dockerfilePath string) string {

	// needs using goroutine
	// and saving log line by line

	go buildJob(repoName, dockerfilePath)

	// only ok
	return `{"message":"ok"}`
}

func buildJob(repoName string, dockerfilePath string) {
	logger.DEBUG("docker.go", "buildJob start")

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

	logger.DEBUG("docker.go", "buildJob end")
}
