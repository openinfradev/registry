package model

// TagResult is tag name and digest, size
type TagResult struct {
	Name   string `json:"name"`
	Digest string `json:"digest"`
	// Size   string `json:"size"`
}

// CatalogResult is registry catalog result
type CatalogResult struct {
	Repositories []string `json:"repositories"`
}

// RepositoryRaw is registry repository raw data
type RepositoryRaw struct {
	Name string   `json:"name"`
	Tags []string `json:"tags"`
}

// RepositoryResult is registry repository result
type RepositoryResult struct {
	Name string      `json:"name"`
	Tags []TagResult `json:"tags"`
}

// RepositoriesResult is regitry repositories result
type RepositoriesResult struct {
	Repositories []RepositoryResult `json:"repositories"`
}

// RegistryCommonCode is test struct
type RegistryCommonCode struct {
	CodeName  string `json:"codeName"`
	GroupCode string `json:"groupCode"`
}
