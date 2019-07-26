package service

import (
	"bufio"
	"fmt"
	"os/exec"
)

// DockerService is docker command relative service
type DockerService struct{}

// Build returns docker building logs with dockerfile
func (d *DockerService) Build(repoName string, dockerfilePath string) string {

	// needs using goroutine
	// and saving log line by line

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
		fmt.Println(m)
	}
	build.Wait()

	return r
}
