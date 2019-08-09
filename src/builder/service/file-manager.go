package service

import (
	"builder/util"
	"builder/util/logger"
	"errors"
	"io/ioutil"
	"os"
)

// FileManager is file manager
type FileManager struct{}

// DeleteDirectory is path removing all
func (f *FileManager) DeleteDirectory(path string) {
	err := os.RemoveAll(path)
	if err != nil {
		logger.ERROR("file-manager.go", err.Error())
	}
}

// WriteDockerfile returns written dockerfile path
func (f *FileManager) WriteDockerfile(contents string) (string, error) {
	dirPath := basicinfo.TemporaryPath + "/" + util.GetTimeMillisecond()
	err := os.Mkdir(dirPath, 0755)
	if err == nil || os.IsExist(err) {
		ioutil.WriteFile(dirPath+"/Dockerfile", []byte(contents), 0755)
		return dirPath, nil
	}
	return "", errors.New("Failed to make temporary directory")
}

// PullGitRepository returns pulled git repository dockerfile path
func (f *FileManager) PullGitRepository(gitRepo string, userID string, userPW string) (string, error) {

	// using go-routine but sync

	return "", nil
}
