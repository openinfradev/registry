package service

import (
	"builder/config"
	"bufio"
	urlconst "builder/constant/url"
	"builder/util"
	"builder/util/logger"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"os/exec"
)

// FileManager is file manager
type FileManager struct{}

// GetTemporaryPath returns temporary path(directory) string
func (f *FileManager) GetTemporaryPath() string {
	return config.GetConfig().Default.TmpDir + "/" + util.GetTimeMillisecond()
}

// MakeDirectory is making dir on root
func (f *FileManager) MakeDirectory(rootDir string, dir string) string {
	target := fmt.Sprintf("%s/%s", rootDir, dir)
	if _, err := os.Stat(target); os.IsNotExist(err) {
		os.Mkdir(target, os.ModeDir)
	}
	return target
}

// DeleteDirectory is path removing all
func (f *FileManager) DeleteDirectory(path string) {
	err := os.RemoveAll(path)
	if err != nil {
		logger.ERROR("service/file-manager.go", "DeleteDirectory", err.Error())
	}
}

// CopyFile is copy file src to dist
func (f *FileManager) CopyFile(source string, dest string) (err error) {
	sourcefile, err := os.Open(source)
	if err != nil {
		return err
	}

	defer sourcefile.Close()

	destfile, err := os.Create(dest)
	if err != nil {
		return err
	}

	defer destfile.Close()

	_, err = io.Copy(destfile, sourcefile)
	if err == nil {
		sourceinfo, err := os.Stat(source)
		if err != nil {
			err = os.Chmod(dest, sourceinfo.Mode())
		}
	}
	return
}

// CopyDirectory is directory src to dist
func (f *FileManager) CopyDirectory(source string, dest string) (err error) {

	// get properties of source dir
	sourceinfo, err := os.Stat(source)
	if err != nil {
		return err
	}

	// create dest dir
	err = os.MkdirAll(dest, sourceinfo.Mode())
	if err != nil {
		return err
	}

	directory, _ := os.Open(source)
	objects, err := directory.Readdir(-1)

	for _, obj := range objects {
		sourcefilepointer := source + "/" + obj.Name()
		destinationfilepointer := dest + "/" + obj.Name()

		if obj.IsDir() {
			// create sub-directories - recursively
			err = f.CopyDirectory(sourcefilepointer, destinationfilepointer)
			if err != nil {
				logger.ERROR("service/file-manager.go", "CopyDirectory", err.Error())
			}
		} else {
			// perform copy
			err = f.CopyFile(sourcefilepointer, destinationfilepointer)
			if err != nil {
				logger.ERROR("service/file-manager.go", "CopyDirectory", err.Error())
			}
		}
	}
	return
}


// WriteDockerfile returns written dockerfile path
func (f *FileManager) WriteDockerfile(contents string) (string, error) {
	dirPath := f.GetTemporaryPath()
	err := os.Mkdir(dirPath, 0755)
	if err == nil || os.IsExist(err) {
		ioutil.WriteFile(dirPath+"/Dockerfile", []byte(contents), 0755)
		return dirPath, nil
	}
	return "", errors.New("Failed to make temporary directory")
}

// PullGitRepository returns pulled git repository dockerfile path
func (f *FileManager) PullGitRepository(gitRepoURL string, userID string, userPW string) (string, error) {
	dirPath := f.GetTemporaryPath()
	gitRepo := util.ExtractGitRepositoryURL(gitRepoURL)
	gitURL := ""
	if userID == "" || userPW == "" {
		// public
		gitURL = fmt.Sprintf(urlconst.GitRepositoryPublicURL, gitRepo.Protocol, gitRepo.URL)
	} else {
		// private
		gitURL = fmt.Sprintf(urlconst.GitRepositoryPrivateURL, gitRepo.Protocol, url.QueryEscape(userID), url.QueryEscape(userPW), gitRepo.URL)
	}
	gitClone := exec.Command("docker", "run", "--rm", "-v", fmt.Sprintf("%s:/tmp", config.GetConfig().Default.TmpDir), "alpine/git", "clone", gitURL, dirPath)

	logger.DEBUG("service/file-manager.go", "PullGitRepository", gitURL)

	// make stdout pipeline but anything doesn't printed.
	stdout, err := gitClone.StdoutPipe()
	if err != nil {
		return "", errors.New("Failed to make temporary directory")
	}
	gitClone.Start()
	scanner := bufio.NewScanner(stdout)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		m := scanner.Text()
		logger.DEBUG("service/file-manager.go", "PullGitRepository", m)
	}
	gitClone.Wait()

	return dirPath, nil
}
