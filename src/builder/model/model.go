package model

// BasicResult is basic result for response
type BasicResult struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// CatalogResult is registry catalog result
type CatalogResult struct {
	Repositories []string `json:"repositories"`
}

// RepositoryResult is registry repository result
type RepositoryResult struct {
	Name string   `json:"name"`
	Tags []string `json:"tags"`
}

// RepositoriesResult is regitry repositories result
type RepositoriesResult struct {
	Repositories []RepositoryResult `json:"repositories"`
}
