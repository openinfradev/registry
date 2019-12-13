package service

import (
	"builder/config"
	"builder/constant/minio"
	"bufio"
	"builder/constant"
	tacoconst "builder/constant/taco"
	"builder/model"
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


// BuildByCopiedMinioBucket is docker buiding by copied minio bucket
func (d *DockerService) BuildByCopiedMinioBucket(params *model.DockerBuildByMinioCopyAsParam) *model.BasicResult {

	// phase - preparing
	p := tacoutil.MakePhaseLog(params.BuildID, tacoconst.PhasePreparing.StartSeq, tacoconst.PhasePreparing.Status)
	is.RegistryRepository.InsertBuildLog(p)

	// phase - unpacking
	p = tacoutil.MakePhaseLog(params.BuildID, tacoconst.PhaseUnpacking.StartSeq, tacoconst.PhaseUnpacking.Status)
	is.RegistryRepository.InsertBuildLog(p)

	// copy as
	src := params.SrcPath
	if strings.HasPrefix(src, "/") {
		src = strings.Replace(src, "/", "", 1)
	}
	dest := params.Path
	if strings.HasPrefix(dest, "/") {
		dest = strings.Replace(dest, "/", "", 1)
	}
	srcPath := fmt.Sprintf("%s/%s/%s", minio.MinioDataPath, params.SrcUserID, src)
	destPath := fmt.Sprintf("%s/%s/%s", minio.MinioDataPath, params.UserID, dest)
	err := is.FileManager.CopyDirectory(srcPath, destPath)
	if err != nil {
		return &model.BasicResult{
			Code: constant.ResultFail,
			Message: "",
		}
	}

	return d.Build(params.BuildID, params.Name, params.Tag, destPath, params.UseCache, params.Push, false)
}

// BuildByMinioBucket is docker building by minio bucket
func (d *DockerService) BuildByMinioBucket(params *model.DockerBuildByMinioParam) *model.BasicResult {
	// phase - preparing
	p := tacoutil.MakePhaseLog(params.BuildID, tacoconst.PhasePreparing.StartSeq, tacoconst.PhasePreparing.Status)
	is.RegistryRepository.InsertBuildLog(p)

	// phase - unpacking
	p = tacoutil.MakePhaseLog(params.BuildID, tacoconst.PhaseUnpacking.StartSeq, tacoconst.PhaseUnpacking.Status)
	is.RegistryRepository.InsertBuildLog(p)

	path := params.Path
	if strings.HasPrefix(path, "/") {
		path = strings.Replace(path, "/", "", 1)
	}
	fullPath := fmt.Sprintf("%s/%s/%s", minio.MinioDataPath, params.UserID, path)
	return d.Build(params.BuildID, params.Name, params.Tag, fullPath, params.UseCache, params.Push, false)
}

// BuildByDockerfile is docker building by dockerfile
func (d *DockerService) BuildByDockerfile(params *model.DockerBuildByFileParam) *model.BasicResult {

	// phase - preparing
	p := tacoutil.MakePhaseLog(params.BuildID, tacoconst.PhasePreparing.StartSeq, tacoconst.PhasePreparing.Status)
	is.RegistryRepository.InsertBuildLog(p)

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
	is.RegistryRepository.InsertBuildLog(p)

	path, err := is.FileManager.WriteDockerfile(string(decoded))
	if err != nil {
		return &model.BasicResult{
			Code:    constant.ResultFail,
			Message: "",
		}
	}

	return d.Build(params.BuildID, params.Name, params.Tag, path, params.UseCache, params.Push, true)
}

// BuildByGitRepository is docker building by git repository
func (d *DockerService) BuildByGitRepository(params *model.DockerBuildByGitParam) *model.BasicResult {

	// phase - preparing
	p := tacoutil.MakePhaseLog(params.BuildID, tacoconst.PhasePreparing.StartSeq, tacoconst.PhasePreparing.Status)
	is.RegistryRepository.InsertBuildLog(p)

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
	is.RegistryRepository.InsertBuildLog(p)

	// not using go-routine (not yet)
	// ch := make(chan string, 1)	// dirPath
	path, err := is.FileManager.PullGitRepository(params.GitRepository, params.UserID, string(decoded))
	if err != nil {
		return &model.BasicResult{
			Code:    constant.ResultFail,
			Message: "",
		}
	}

	return d.Build(params.BuildID, params.Name, params.Tag, path, params.UseCache, params.Push, true)
}

// Build is docker building by file path
func (d *DockerService) Build(buildID string, repoName string, tag string, dockerfilePath string, useCache bool, push bool, tempDelete bool) *model.BasicResult {

	// async
	ch := make(chan string)
	if push {
		go d.BuildAndPush(ch, buildID, repoName, tag, dockerfilePath, useCache, tempDelete)
	} else {
		go buildJob(ch, buildID, repoName, tag, dockerfilePath, useCache, tempDelete)
	}

	// only ok
	return &model.BasicResult{
		Code:    constant.ResultSuccess,
		Message: "",
	}
}

// BuildAndPush is docker build and push
func (d *DockerService) BuildAndPush(ch chan<- string, buildID string, repoName string, tag string, dockerfilePath string, useCache bool, tempDelete bool) {
	proc := make(chan string)
	// build
	go buildJob(proc, buildID, repoName, tag, dockerfilePath, useCache, tempDelete)
	r := <-proc
	if r == constant.ResultFail {
		procBuildError(buildID, repoName, tag)
		ch <- constant.ResultFail
		return
	}

	// push
	// phase - push
	is.RegistryRepository.UpdateBuildPhase(buildID, tacoconst.PhasePushing.Status)
	p := tacoutil.MakePhaseLog(buildID, tacoconst.PhasePushing.StartSeq, tacoconst.PhasePushing.Status)
	is.RegistryRepository.InsertBuildLog(p)

	go pushJob(proc, repoName, tag)
	r = <-proc
	if r == constant.ResultFail {
		procBuildError(buildID, repoName, tag)
		ch <- constant.ResultFail
		return
	}

	// security scan (optional??)
	// returned value isn't necessary.
	is.SecurityService.Scan(repoName, tag)

	// if not exists "latest" tag. latest tagging and push
	latest := "latest"
	if !is.RegistryService.ExistsTag(repoName, latest) {
		tagParams := &model.DockerTagParam{
			BuildID : buildID,
			Name : repoName,
			OldTag : tag,
			NewTag : latest,
		}
		go d.PullAndTag(proc, tagParams)
		r = <-proc
		if r == constant.ResultFail {
			logger.ERROR("service/docker.go", "BuildAndPush", fmt.Sprintf("Failed to tagging 'latest' on [%s].", repoName))
		}
	}

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
		procTagError(params.BuildID, params.Name, params.NewTag)
		ch <- constant.ResultFail
		return
	}

	// 2. tag
	go tagJob(proc, params.Name, params.OldTag, params.NewTag)
	r = <-proc
	if r == constant.ResultFail {
		logger.ERROR("service/docker.go", "PullAndTag", "failed to tagging docker image")
		procTagError(params.BuildID, params.Name, params.NewTag)
		ch <- constant.ResultFail
		return
	}

	// 3. push
	go pushJob(proc, params.Name, params.NewTag)
	r = <-proc
	if r == constant.ResultFail {
		logger.ERROR("service/docker.go", "PullAndTag", "failed to pushing docker image")
		procTagError(params.BuildID, params.Name, params.NewTag)
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

	registryinfo := config.GetConfig().Registry

	login := exec.Command("docker", "login", registryinfo.Endpoint, "--username", tacoconst.BuilderUser, "--password", tacoconst.BuilderPass)
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
		logger.INFO("service/docker.go", "loginJob", fmt.Sprintf("[%s] logged in succeeded", registryinfo.Endpoint))
		ch <- constant.ResultSuccess
	} else {
		logger.INFO("service/docker.go", "loginJob", fmt.Sprintf("[%s] logged in failed", registryinfo.Endpoint))
		ch <- constant.ResultFail
	}
}

func pullJob(ch chan<- string, repoName string, tag string, external bool) {

	registryinfo := config.GetConfig().Registry

	logger.DEBUG("service/docker.go", "pullJob", fmt.Sprintf("pullJob start [%s:%s]", repoName, tag))

	if external {
		repoName = fmt.Sprintf("%s:%s", repoName, tag)
	} else {
		repoName = fmt.Sprintf("%s/%s:%s", registryinfo.Endpoint, repoName, tag)
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

	registryinfo := config.GetConfig().Registry

	logger.DEBUG("service/docker.go", "pushJob", fmt.Sprintf("pushJob start [%s:%s]", repoName, tag))

	repoName = fmt.Sprintf("%s/%s:%s", registryinfo.Endpoint, repoName, tag)
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

	registryinfo := config.GetConfig().Registry

	logger.DEBUG("service/docker.go", "tagJob", fmt.Sprintf("tagJob [%s] [%s] to [%s]", repoName, oldTag, newTag))

	oldRepo := fmt.Sprintf("%s/%s:%s", registryinfo.Endpoint, repoName, oldTag)
	newRepo := fmt.Sprintf("%s/%s:%s", registryinfo.Endpoint, repoName, newTag)

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

func buildJob(ch chan<- string, buildID string, repoName string, tag string, dockerfilePath string, useCache bool, tempDelete bool) {

	registryinfo := config.GetConfig().Registry

	logger.DEBUG("service/docker.go", "buildJob", "buildJob start "+repoName)

	seq := tacoconst.PhaseBuilding.StartSeq

	// phase - build
	// updating build phase
	is.RegistryRepository.UpdateBuildPhase(buildID, tacoconst.PhaseBuilding.Status)
	p := tacoutil.MakePhaseLog(buildID, seq, tacoconst.PhaseBuilding.Status)
	is.RegistryRepository.InsertBuildLog(p)

	repoName = fmt.Sprintf("%s/%s:%s", registryinfo.Endpoint, repoName, tag)
	var build *exec.Cmd
	if useCache {
		// phase - checking cache
		p = tacoutil.MakePhaseLog(buildID, seq, tacoconst.PhaseCheckingCache.Status)
		is.RegistryRepository.InsertBuildLog(p)

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
		is.RegistryRepository.InsertBuildLog(row)

		// logger.DEBUG("service/docker.go", "buildJob", m)
	}
	errscan := bufio.NewScanner(stderr)
	errscan.Split(bufio.ScanLines)
	for errscan.Scan() {
		seq++
		m := errscan.Text()
		errrow := tacoutil.ParseLog(buildID, seq, m)
		errrow.Type = "error"
		is.RegistryRepository.InsertBuildLog(errrow)
		logger.ERROR("service/docker.go", "buildJob", m)
		// path removeall - because of breaking
		if tempDelete {
			is.FileManager.DeleteDirectory(dockerfilePath)
		}

		ch <- constant.ResultFail
	}

	build.Wait()

	// path removeall
	if tempDelete {
		is.FileManager.DeleteDirectory(dockerfilePath)
	}

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
	digest := is.RegistryService.GetDigest(repoName, tag)
	size := getImageSize(repoName, tag)
	is.RegistryRepository.UpdateTagDigest(buildID, tag, digest, size)

	is.RegistryRepository.UpdateBuildPhase(buildID, tacoconst.PhaseComplete.Status)
	p := tacoutil.MakePhaseLog(buildID, tacoconst.PhaseComplete.StartSeq, tacoconst.PhaseComplete.Status)
	is.RegistryRepository.InsertBuildLog(p)
}

func procBuildError(buildID string, repoName string, tag string) {
	is.RegistryRepository.DeleteUsageLog(buildID, tag)
	if !is.RegistryService.ExistsTag(repoName, tag) {
		is.RegistryRepository.DeleteTag(buildID, tag)
	}

	is.RegistryRepository.UpdateBuildPhase(buildID, tacoconst.PhaseError.Status)
	p := tacoutil.MakePhaseLog(buildID, tacoconst.PhaseError.StartSeq, tacoconst.PhaseError.Status)
	is.RegistryRepository.InsertBuildLog(p)
}

func procTagComplete(buildID string, repoName string, tag string) {
	// digest & size
	digest := is.RegistryService.GetDigest(repoName, tag)
	size := getImageSize(repoName, tag)
	is.RegistryRepository.UpdateTagDigest(buildID, tag, digest, size)
}

func procTagError(buildID string, repoName string, tag string) {
	if !is.RegistryService.ExistsTag(repoName, tag) {
		is.RegistryRepository.DeleteTag(buildID, tag)
	}
}

func getImageSize(repoName string, tag string) string {

	registryinfo := config.GetConfig().Registry

	repo := fmt.Sprintf("%s/%s:%s", registryinfo.Endpoint, repoName, tag)
	cmd := "docker images --filter=reference='" + repo + "' --format \"{{.Size}}\""
	stdout, err := exec.Command("/bin/sh", "-c", cmd).Output()
	if err != nil {
		logger.ERROR("service/docker.go", "getImageSize", err.Error())
		return "0"
	}
	return string(stdout)
}
