package service

import (
	"builder/constant/minio"
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
	"strings"
	"time"
)

// DockerService is docker command relative service
type DockerService struct{}

var fileManager *FileManager
var registryRepository *repository.RegistryRepository
var registryService *RegistryService
var securityService *SecurityService

func init() {
	fileManager = new(FileManager)
	registryRepository = new(repository.RegistryRepository)
	registryService = new(RegistryService)
	securityService = new(SecurityService)
}

// BuildByCopiedMinioBucket is docker buiding by copied minio bucket
func (d *DockerService) BuildByCopiedMinioBucket() *model.BasicResult {
	return nil
}

// BuildByMinioBucket is docker building by minio bucket
func (d *DockerService) BuildByMinioBucket(params *model.DockerBuildByMinioParam) *model.BasicResult {
	// phase - preparing
	p := tacoutil.MakePhaseLog(params.BuildID, tacoconst.PhasePreparing.StartSeq, tacoconst.PhasePreparing.Status)
	registryRepository.InsertBuildLog(p)

	// phase - unpacking
	p = tacoutil.MakePhaseLog(params.BuildID, tacoconst.PhaseUnpacking.StartSeq, tacoconst.PhaseUnpacking.Status)
	registryRepository.InsertBuildLog(p)

	path := fmt.Sprintf("%s/%s", minio.MinioDataPath, params.UserID)
	return d.Build(params.BuildID, params.Name, path, params.UseCache, params.Push)
}

// BuildByDockerfile is docker building by dockerfile
func (d *DockerService) BuildByDockerfile(params *model.DockerBuildByFileParam) *model.BasicResult {

	// phase - preparing
	p := tacoutil.MakePhaseLog(params.BuildID, tacoconst.PhasePreparing.StartSeq, tacoconst.PhasePreparing.Status)
	registryRepository.InsertBuildLog(p)

	// decoding contents
	decoded, err := base64.StdEncoding.DecodeString(params.Contents)
	if err != nil {
		return &model.BasicResult{
			Code:    constant.ResultFail,
			Message: "contents isn't base64 encoded",
		}
	}

	// phase - unpacking
	p = tacoutil.MakePhaseLog(params.BuildID, tacoconst.PhaseUnpacking.StartSeq, tacoconst.PhaseUnpacking.Status)
	registryRepository.InsertBuildLog(p)

	path, err := fileManager.WriteDockerfile(string(decoded))
	if err != nil {
		return &model.BasicResult{
			Code:    constant.ResultFail,
			Message: "",
		}
	}

	return d.Build(params.BuildID, params.Name, path, params.UseCache, params.Push)
}

// BuildByGitRepository is docker building by git repository
func (d *DockerService) BuildByGitRepository(params *model.DockerBuildByGitParam) *model.BasicResult {

	// phase - preparing
	p := tacoutil.MakePhaseLog(params.BuildID, tacoconst.PhasePreparing.StartSeq, tacoconst.PhasePreparing.Status)
	registryRepository.InsertBuildLog(p)

	// decoding userPW
	decoded, err := base64.StdEncoding.DecodeString(params.UserPW)
	if err != nil {
		return &model.BasicResult{
			Code:    constant.ResultFail,
			Message: "userPW isn't base64 encoded",
		}
	}

	// phase - unpacking
	p = tacoutil.MakePhaseLog(params.BuildID, tacoconst.PhaseUnpacking.StartSeq, tacoconst.PhaseUnpacking.Status)
	registryRepository.InsertBuildLog(p)

	// not using go-routine (not yet)
	// ch := make(chan string, 1)	// dirPath
	path, err := fileManager.PullGitRepository(params.GitRepository, params.UserID, string(decoded))
	if err != nil {
		return &model.BasicResult{
			Code:    constant.ResultFail,
			Message: "",
		}
	}

	return d.Build(params.BuildID, params.Name, path, params.UseCache, params.Push)
}

// Build is docker building by file path
func (d *DockerService) Build(buildID string, repoName string, dockerfilePath string, useCache bool, push bool) *model.BasicResult {

	// async
	ch := make(chan string)
	if push {
		go d.BuildAndPush(ch, buildID, repoName, dockerfilePath, useCache)
	} else {
		go buildJob(ch, buildID, repoName, dockerfilePath, useCache)
	}

	// only ok
	return &model.BasicResult{
		Code:    constant.ResultSuccess,
		Message: "",
	}
}

// BuildAndPush is docker build and push
func (d *DockerService) BuildAndPush(ch chan<- string, buildID string, repoName string, dockerfilePath string, useCache bool) {
	// fixed "latest"
	tag := "latest"

	proc := make(chan string)
	// build
	go buildJob(proc, buildID, repoName, dockerfilePath, useCache)
	r := <-proc
	if r == constant.ResultFail {
		procBuildError(buildID)
		ch <- constant.ResultFail
		return
	}

	// tag
	// go tagJob(proc, repoName, tag, tag)
	// r = <-proc
	// if r == constant.ResultFail {
	// 	procBuildError(buildID)
	// 	ch <- constant.ResultFail
	// 	return
	// }

	// push
	// phase - push
	registryRepository.UpdateBuildPhase(buildID, tacoconst.PhasePushing.Status)
	p := tacoutil.MakePhaseLog(buildID, tacoconst.PhasePushing.StartSeq, tacoconst.PhasePushing.Status)
	registryRepository.InsertBuildLog(p)

	go pushJob(proc, repoName, tag)
	r = <-proc
	if r == constant.ResultFail {
		procBuildError(buildID)
		ch <- constant.ResultFail
		return
	}

	// security scan (optional??)
	// returned value isn't necessary.
	securityService.Scan(repoName, tag)

	// phase - complete
	procBuildComplete(buildID, repoName, tag)
	ch <- constant.ResultSuccess
}

// PullAndTag is docker image pulling and tagging
func (d *DockerService) PullAndTag(ch chan<- string, params *model.DockerTagParam) {
	proc := make(chan string)
	logger.DEBUG("service/docker.go", "PullAndTag", fmt.Sprintf("start %s from %s to %s", params.Name, params.OldTag, params.NewTag))

	// 1. pull
	go pullJob(proc, params.Name, params.OldTag, false)
	r := <-proc
	if r == constant.ResultFail {
		logger.ERROR("service/docker.go", "PullAndTag", "failed to pulling docker image")
		procTagError(params.BuildID, params.NewTag)
		ch <- constant.ResultFail
		return
	}

	// 2. tag
	go tagJob(proc, params.Name, params.OldTag, params.NewTag)
	r = <-proc
	if r == constant.ResultFail {
		logger.ERROR("service/docker.go", "PullAndTag", "failed to tagging docker image")
		procTagError(params.BuildID, params.NewTag)
		ch <- constant.ResultFail
		return
	}

	// 3. push
	go pushJob(proc, params.Name, params.NewTag)
	r = <-proc
	if r == constant.ResultFail {
		logger.ERROR("service/docker.go", "PullAndTag", "failed to pushing docker image")
		procTagError(params.BuildID, params.NewTag)
		ch <- constant.ResultFail
		return
	}

	logger.DEBUG("service/docker.go", "PullAndTag", fmt.Sprintf("end %s from %s to %s", params.Name, params.OldTag, params.NewTag))
	procTagComplete(params.BuildID, params.Name, params.NewTag)
	ch <- constant.ResultSuccess
}

// Tag is image tagging
func (d *DockerService) Tag(params *model.DockerTagParam) *model.BasicResult {

	// sync
	ch := make(chan string)
	go d.PullAndTag(ch, params)
	r := <-ch

	return &model.BasicResult{
		Code:    r,
		Message: "",
	}
}

// Push is docker image pushing
func (d *DockerService) Push(params *model.DockerPushParam) *model.BasicResult {

	// async
	ch := make(chan string)
	go pushJob(ch, params.Name, params.Tag)

	// only ok
	return &model.BasicResult{
		Code:    constant.ResultSuccess,
		Message: "",
	}
}

// Pull is docker image pulling
func (d *DockerService) Pull(params *model.DockerPullParam, async bool, external bool) *model.BasicResult {

	// async
	ch := make(chan string)
	go pullJob(ch, params.Name, params.Tag, external)
	r := constant.ResultSuccess
	if !async {
		r = <-ch
	}

	return &model.BasicResult{
		Code:    r,
		Message: "",
	}
}

// Login is registry logged in
func (d *DockerService) Login() {

	ch := make(chan string)
	ticker := time.NewTicker(time.Second * 5)
	for t := range ticker.C {
		logger.DEBUG("service/docker.go", "Login", "try "+t.String())
		go loginJob(ch)
		r := <-ch
		if r == constant.ResultSuccess {
			logger.DEBUG("service/docker.go", "Login", "complete")
			ticker.Stop()
		}
	}
}

func loginJob(ch chan<- string) {

	login := exec.Command("docker", "login", basicinfo.RegistryEndpoint, "--username", tacoconst.BuilderUser, "--password", tacoconst.BuilderPass)
	r := ""
	stdout, _ := login.StdoutPipe()
	login.Start()
	scanner := bufio.NewScanner(stdout)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		m := scanner.Text()
		r += m + "\n"
	}
	login.Wait()

	if strings.Contains(r, tacoconst.LoginSucceeded) {
		logger.INFO("service/docker.go", "loginJob", fmt.Sprintf("[%s] logged in succeeded", basicinfo.RegistryEndpoint))
		ch <- constant.ResultSuccess
	} else {
		logger.INFO("service/docker.go", "loginJob", fmt.Sprintf("[%s] logged in failed", basicinfo.RegistryEndpoint))
		ch <- constant.ResultFail
	}
}

func pullJob(ch chan<- string, repoName string, tag string, external bool) {
	logger.DEBUG("service/docker.go", "pullJob", fmt.Sprintf("pullJob start [%s:%s]", repoName, tag))

	if external {
		repoName = repoName + ":" + tag
	} else {
		repoName = basicinfo.RegistryEndpoint + "/" + repoName + ":" + tag
	}

	pull := exec.Command("docker", "pull", repoName)

	r := ""
	stdout, _ := pull.StdoutPipe()
	stderr, _ := pull.StderrPipe()
	pull.Start()
	scanner := bufio.NewScanner(stdout)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		m := scanner.Text()
		r += m + "\n"
		// logger.DEBUG("service/docker.go", "pullJob", m)
	}
	errscan := bufio.NewScanner(stderr)
	errscan.Split(bufio.ScanLines)
	for errscan.Scan() {
		m := errscan.Text()
		logger.ERROR("service/docker.go", "pullJob", m)

		ch <- constant.ResultFail
	}
	pull.Wait()

	ch <- constant.ResultSuccess

	logger.DEBUG("service/docker.go", "pullJob", fmt.Sprintf("pullJob end [%s]", repoName))
}

func pushJob(ch chan<- string, repoName string, tag string) {
	logger.DEBUG("service/docker.go", "pushJob", fmt.Sprintf("pushJob start [%s:%s]", repoName, tag))

	repoName = basicinfo.RegistryEndpoint + "/" + repoName + ":" + tag
	push := exec.Command("docker", "push", repoName)

	r := ""
	stdout, _ := push.StdoutPipe()
	stderr, _ := push.StderrPipe()
	push.Start()
	scanner := bufio.NewScanner(stdout)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		m := scanner.Text()
		r += m + "\n"
		// logger.DEBUG("service/docker.go", "pushJob", m)
	}
	errscan := bufio.NewScanner(stderr)
	errscan.Split(bufio.ScanLines)
	for errscan.Scan() {
		m := errscan.Text()
		logger.ERROR("service/docker.go", "pushJob", m)

		ch <- constant.ResultFail
	}
	push.Wait()

	ch <- constant.ResultSuccess

	logger.DEBUG("service/docker.go", "pushJob", fmt.Sprintf("pushJob end [%s]", repoName))
}

func tagJob(ch chan<- string, repoName string, oldTag string, newTag string) {
	logger.DEBUG("service/docker.go", "tagJob", fmt.Sprintf("tagJob [%s] [%s] to [%s]", repoName, oldTag, newTag))

	oldRepo := basicinfo.RegistryEndpoint + "/" + repoName + ":" + oldTag
	newRepo := basicinfo.RegistryEndpoint + "/" + repoName + ":" + newTag

	tag := exec.Command("docker", "tag", oldRepo, newRepo)

	err := tag.Run()
	if err != nil {
		logger.ERROR("service/docker.go", "tagJob", "tagJob is failed")
		ch <- constant.ResultFail
	} else {
		logger.DEBUG("service/docker.go", "tagJob", "tagJob is success")
		ch <- constant.ResultSuccess
	}
}

func buildJob(ch chan<- string, buildID string, repoName string, dockerfilePath string, useCache bool) {
	logger.DEBUG("service/docker.go", "buildJob", "buildJob start "+repoName)

	seq := tacoconst.PhaseBuilding.StartSeq

	// phase - build
	// updating build phase
	registryRepository.UpdateBuildPhase(buildID, tacoconst.PhaseBuilding.Status)
	p := tacoutil.MakePhaseLog(buildID, seq, tacoconst.PhaseBuilding.Status)
	registryRepository.InsertBuildLog(p)

	repoName = basicinfo.RegistryEndpoint + "/" + repoName + ":latest"
	var build *exec.Cmd
	if useCache {
		// phase - checking cache
		p = tacoutil.MakePhaseLog(buildID, seq, tacoconst.PhaseCheckingCache.Status)
		registryRepository.InsertBuildLog(p)

		build = exec.Command("docker", "build", "--network=host", "-t", repoName, dockerfilePath)
	} else {
		build = exec.Command("docker", "build", "--no-cache", "--network=host", "-t", repoName, dockerfilePath)
	}

	stdout, _ := build.StdoutPipe()
	stderr, _ := build.StderrPipe()
	build.Start()
	scanner := bufio.NewScanner(stdout)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		seq++
		m := scanner.Text()
		row := tacoutil.ParseLog(buildID, seq, m)
		registryRepository.InsertBuildLog(row)

		// logger.DEBUG("service/docker.go", "buildJob", m)
	}
	errscan := bufio.NewScanner(stderr)
	errscan.Split(bufio.ScanLines)
	for errscan.Scan() {
		seq++
		m := errscan.Text()
		errrow := tacoutil.ParseLog(buildID, seq, m)
		errrow.Type = "error"
		registryRepository.InsertBuildLog(errrow)
		logger.ERROR("service/docker.go", "buildJob", m)
		// path removeall - because of breaking
		fileManager.DeleteDirectory(dockerfilePath)

		ch <- constant.ResultFail
	}

	build.Wait()

	// path removeall
	fileManager.DeleteDirectory(dockerfilePath)

	ch <- constant.ResultSuccess

	logger.DEBUG("service/docker.go", "buildJob", "buildJob end "+repoName)
}

// deprecated : has a problem
// func garbageCollectJob(ch chan<- string) {
// 	logger.DEBUG("service/docker.go", "garbageCollectJob", "garbage collect start")

// 	gc := exec.Command("docker", "exec", basicinfo.RegistryName, "bin/registry", "garbage-collect", "/etc/docker/registry/config.yml")

// 	r := ""
// 	stdout, _ := gc.StdoutPipe()
// 	gc.Start()
// 	scanner := bufio.NewScanner(stdout)
// 	scanner.Split(bufio.ScanLines)
// 	for scanner.Scan() {
// 		m := scanner.Text()
// 		r += m + "\n"
// 		// logger.DEBUG("service/docker.go", "garbageCollectJob", m)
// 	}
// 	gc.Wait()

// 	logger.DEBUG("service/docker.go", "garbageCollectJob", "garbage collect end")

// 	ch <- r
// }

func procBuildComplete(buildID string, repoName string, tag string) {
	// digest & size
	digest := registryService.GetDigest(repoName, tag)
	size := getImageSize(repoName, tag)
	registryRepository.UpdateTagDigest(buildID, "latest", digest, size)

	registryRepository.UpdateBuildPhase(buildID, tacoconst.PhaseComplete.Status)
	p := tacoutil.MakePhaseLog(buildID, tacoconst.PhaseComplete.StartSeq, tacoconst.PhaseComplete.Status)
	registryRepository.InsertBuildLog(p)
}

func procBuildError(buildID string) {
	registryRepository.DeleteUsageLog(buildID)
	registryRepository.DeleteTag(buildID, "latest")

	registryRepository.UpdateBuildPhase(buildID, tacoconst.PhaseError.Status)
	p := tacoutil.MakePhaseLog(buildID, tacoconst.PhaseError.StartSeq, tacoconst.PhaseError.Status)
	registryRepository.InsertBuildLog(p)
}

func procTagComplete(buildID string, repoName string, tag string) {
	// digest & size
	digest := registryService.GetDigest(repoName, tag)
	size := getImageSize(repoName, tag)
	registryRepository.UpdateTagDigest(buildID, tag, digest, size)
}

func procTagError(buildID string, tag string) {
	registryRepository.DeleteTag(buildID, tag)
}

func getImageSize(repoName string, tag string) string {

	repo := basicinfo.RegistryEndpoint + "/" + repoName + ":" + tag
	cmd := "docker images --filter=reference='" + repo + "' --format \"{{.Size}}\""
	stdout, err := exec.Command("/bin/sh", "-c", cmd).Output()
	if err != nil {
		logger.ERROR("service/docker.go", "getImageSize", err.Error())
		return "0"
	}
	return string(stdout)
}
