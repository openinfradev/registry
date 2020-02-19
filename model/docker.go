package model

// DockerBuildParam is base parameters for docker build
type DockerBuildParam struct {
	BuildID  string `json:"build" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Tag      string `json:"tag" binding:"required"`
	UseCache bool   `json:"useCache" binding:"required"`
	Push     bool   `json:"push" binding:"required"`
}

// DockerBuildByFileParam is parameters for docker api
type DockerBuildByFileParam struct {
	DockerBuildParam
	Contents string `json:"contents" binding:"required"`
}

// DockerBuildByGitParam is parameters for docker api
type DockerBuildByGitParam struct {
	DockerBuildParam
	GitRepository string `json:"gitRepo" binding:"required"`
	UserID        string `json:"userId" binding:"required"`
	UserPW        string `json:"userPw" binding:"required"`
}

// DockerBuildByMinioParam is parameters for docker api
type DockerBuildByMinioParam struct {
	DockerBuildParam
	UserID string `json:"userId" binding:"required"`
	Path   string `json:"path" binding:"required"`
}

// DockerBuildByMinioCopyAsParam is parameters for docker api
type DockerBuildByMinioCopyAsParam struct {
	DockerBuildByMinioParam
	SrcUserID string `json:"srcUserId" binding:"required"`
	SrcPath   string `json:"srcPath" binding:"required"`
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

// DockerPullParam is parameters for docker api
type DockerPullParam struct {
	Name string `json:"name" binding:"required"`
	Tag  string `json:"tag" binding:"required"`
}
