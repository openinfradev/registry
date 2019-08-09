package model

// DockerBuildByFileParam is parameters for docker api
type DockerBuildByFileParam struct {
	Name     string `json:"name" binding:"required"`
	Contents string `json:"contents" binding:"required"`
}

// DockerBuildByGitParam is parameters for docker api
type DockerBuildByGitParam struct {
	Name          string `json:"name" binding:"required"`
	GitRepository string `json:"gitRepo" binding:"required"`
	UserID        string `json:"userId" binding:"required"`
	UserPW        string `json:"userPw" binding:"required"`
}

// DockerTagParam is parameters for docker api
type DockerTagParam struct {
	Name   string `json:"name" binding:"required"`
	OldTag string `json:"oldTag" binding:"required"`
	NewTag string `json:"newTag" binding:"required"`
}

// DockerPushParam is parameters for docker api
type DockerPushParam struct {
	Name string `json:"name" binding:"required"`
	Tag  string `json:"tag" binding:"required"`
}
