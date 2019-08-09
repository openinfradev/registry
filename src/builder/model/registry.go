package model

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
