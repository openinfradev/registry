package model

// MinioParam is minio parameter
type MinioParam struct {
	UserID string `json:"userId"`
	UserPW string `json:"userPw"`
}

// MinioResult is minio container domain and port
type MinioResult struct {
	BasicResult
	Domain string `json:"domain"`
	Port   int    `json:"port"`
}
