package service

import (
	"builder/util/logger"
	"bytes"
	"fmt"
	"os/exec"
)

// DockerService is docker command relative service
type DockerService struct{}

// Build returns docker building logs with dockerfile
func (d *DockerService) Build(repoName string, dockerfilePath string) string {
	// using goroutine !!!

	// tag is fixed latest
	// build := "docker build -t " + repoName + ":latest " + dockerfilePath
	// out, err := exec.Command("/bin/sh", "-c", build).Output()
	repoName = repoName + ":latest"
	build := exec.Command("docker", "build", "-t", repoName, dockerfilePath)
	var stdout, stderr bytes.Buffer
	build.Stdout = &stdout
	build.Stderr = &stderr
	err := build.Run()
	fmt.Printf(stdout.String(), stderr.String())
	if err != nil {
		logger.ERROR("service/docker.go", "Failed to build")
		return ""
	}
	// fmt.Printf(stdout.String(), stderr.String())
	return ""
}
