package service

import (
	"bufio"
	"builder/constant"
	tacoconst "builder/constant/taco"
	"builder/model"
	"builder/repository"
	"builder/util/logger"
	tacoutil "builder/util/taco"
	"encoding/base64"
	"fmt"
	"os/exec"
)

// DockerService is docker command relative service
type DockerService struct{}

var fileManager *FileManager
var registryRepository *repository.RegistryRepository

func init() {
	fileManager = new(FileManager)
	registryRepository = new(repository.RegistryRepository)
}

// BuildByDockerfile is docker building by dockerfile
func (d *DockerService) BuildByDockerfile(params *model.DockerBuildByFileParam) *model.BasicResult {
	// needs using goroutine
	// and saving log line by line

	// phase - pulling
	p := tacoutil.MakePhaseLog(params.BuildID, tacoconst.PhasePulling.StartSeq, tacoconst.PhasePulling.Status)
	registryRepository.InsertBuildLog(p)

	// decoding contents
	decoded, err := base64.StdEncoding.DecodeString(params.Contents)
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

	return d.Build(params.BuildID, params.Name, path)
}

// BuildByGitRepository is docker building by git repository
func (d *DockerService) BuildByGitRepository(params *model.DockerBuildByGitParam) *model.BasicResult {
	// needs using goroutine
	// and saving log line by line

	// phase - pulling
	p := tacoutil.MakePhaseLog(params.BuildID, tacoconst.PhasePulling.StartSeq, tacoconst.PhasePulling.Status)
	registryRepository.InsertBuildLog(p)

	// decoding userPW
	decoded, err := base64.StdEncoding.DecodeString(params.UserPW)
	if err != nil {
		return &model.BasicResult{
			Code:    constant.ResultFail,
			Message: "userPW isn't base64 encoded",
		}
	}

	// not using go-routine (not yet)
	// ch := make(chan string, 1)	// dirPath
	path, err := fileManager.PullGitRepository(params.GitRepository, params.UserID, string(decoded))
	if err != nil {
		return &model.BasicResult{
			Code:    constant.ResultFail,
			Message: "",
		}
	}

	return d.Build(params.BuildID, params.Name, path)
}

// Build is docker building by file path
func (d *DockerService) Build(buildID string, repoName string, dockerfilePath string) *model.BasicResult {

	// async
	go buildJob(buildID, repoName, dockerfilePath)

	// only ok
	return &model.BasicResult{
		Code:    constant.ResultSuccess,
		Message: "",
	}
}

// Tag is image tagging
func (d *DockerService) Tag(params *model.DockerTagParam) *model.BasicResult {

	// needs using goroutine
	// and saving log line by line

	// sync
	ch := make(chan model.BasicResult, 1)
	go tagJob(ch, params.Name, params.OldTag, params.NewTag)
	r := <-ch

	return &r
}

// Push is docker image pushing
func (d *DockerService) Push(params *model.DockerPushParam) *model.BasicResult {
	// needs using goroutine
	// and saving log line by line

	// async
	go pushJob(params.Name, params.Tag)

	// only ok
	return &model.BasicResult{
		Code:    constant.ResultSuccess,
		Message: "",
	}
}

func pushJob(repoName string, tag string) {
	logger.DEBUG("service/docker.go", "pushJob", fmt.Sprintf("pushJob start [%s:%s]", repoName, tag))

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
		logger.DEBUG("service/docker.go", "pushJob", m)
	}
	push.Wait()

	logger.DEBUG("service/docker.go", "pushJob", fmt.Sprintf("pushJob end [%s]", repoName))
}

func tagJob(ch chan<- model.BasicResult, repoName string, oldTag string, newTag string) {
	logger.DEBUG("service/docker.go", "tagJob", fmt.Sprintf("tagJob [%s] [%s] to [%s]", repoName, oldTag, newTag))

	result := &model.BasicResult{}

	oldRepo := repoName + ":" + oldTag
	newRepo := basicinfo.RegistryEndpoint + "/" + repoName + ":" + newTag

	tag := exec.Command("docker", "tag", oldRepo, newRepo)

	err := tag.Run()
	if err != nil {
		logger.ERROR("service/docker.go", "tagJob", "tagJob is failed")
		result.Code = constant.ResultFail
		result.Message = ""
		ch <- *result
	} else {
		logger.DEBUG("service/docker.go", "tagJob", "tagJob is success")
		result.Code = constant.ResultSuccess
		result.Message = ""
		ch <- *result
	}
}

func buildJob(buildID string, repoName string, dockerfilePath string) {
	logger.DEBUG("service/docker.go", "buildJob", "buildJob start "+repoName)

	seq := tacoconst.PhaseBuilding.StartSeq

	// phase - build
	p := tacoutil.MakePhaseLog(buildID, seq, tacoconst.PhaseBuilding.Status)
	registryRepository.InsertBuildLog(p)

	repoName = repoName + ":latest"
	build := exec.Command("docker", "build", "--no-cache", "--network=host", "-t", repoName, dockerfilePath)

	rows := []model.BuildLogRow{}
	stdout, _ := build.StdoutPipe()
	build.Start()
	scanner := bufio.NewScanner(stdout)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		seq++
		m := scanner.Text()
		row := tacoutil.ParseLog(buildID, seq, m)
		rows = append(rows, *row)
		logger.DEBUG("service/docker.go", "buildJob", m)
	}
	build.Wait()

	if len(rows) > 0 {
		registryRepository.InsertBuildLogBatch(rows)
	}

	logger.DEBUG("service/docker.go", "buildJob", "buildJob end "+repoName)

	// path removeall
	fileManager.DeleteDirectory(dockerfilePath)

	// phase - complete
	p = tacoutil.MakePhaseLog(buildID, tacoconst.PhaseComplete.StartSeq, tacoconst.PhaseComplete.Status)
	registryRepository.InsertBuildLog(p)
}

func garbageCollectJob(ch chan<- string) {
	logger.DEBUG("service/docker.go", "garbageCollectJob", "garbage collect start")

	gc := exec.Command("docker", "exec", basicinfo.RegistryName, "bin/registry", "garbage-collect", "/etc/docker/registry/config.yml")

	r := ""
	stdout, _ := gc.StdoutPipe()
	gc.Start()
	scanner := bufio.NewScanner(stdout)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		m := scanner.Text()
		r += m + "\n"
		logger.DEBUG("service/docker.go", "garbageCollectJob", m)
	}
	gc.Wait()

	logger.DEBUG("service/docker.go", "garbageCollectJob", "garbage collect end")

	ch <- r
}
