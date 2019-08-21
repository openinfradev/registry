package model

// DockerBuildByFileParam is parameters for docker api
type DockerBuildByFileParam struct {
	BuildID  string `json:"build" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Contents string `json:"contents" binding:"required"`
	UseCache bool   `json:"useCache"`
}

// DockerBuildByGitParam is parameters for docker api
type DockerBuildByGitParam struct {
	BuildID       string `json:"build" binding:"required"`
	Name          string `json:"name" binding:"required"`
	GitRepository string `json:"gitRepo" binding:"required"`
	UserID        string `json:"userId" binding:"required"`
	UserPW        string `json:"userPw" binding:"required"`
	UseCache      bool   `json:"useCache"`
}

// DockerTagParam is parameters for docker api
type DockerTagParam struct {
	BuildID string `json:"build" binding:"required"`
	Name    string `json:"name" binding:"required"`
	OldTag  string `json:"oldTag" binding:"required"`
	NewTag  string `json:"newTag" binding:"required"`
}

// DockerPushParam is parameters for docker api
type DockerPushParam struct {
	BuildID string `json:"build" binding:"required"`
	Name    string `json:"name" binding:"required"`
	Tag     string `json:"tag" binding:"required"`
}
