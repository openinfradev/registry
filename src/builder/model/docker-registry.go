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

// RegistryManifestV1Layer is item
type RegistryManifestV1Layer struct {
	BlobSum string `json:"blobSum"`
}

// RegistryManifestV1History is item
type RegistryManifestV1History struct {
	V1Compatibility string `json:"v1Compatibility"`
}

// RegistryManifestV1Signature is item
type RegistryManifestV1Signature struct {
	// Header    string `json:"header"`
	Signature string `json:"signature"`
	Protected string `json:"protected"`
}

// RegistryManifestV1 is v1 manifest spec
type RegistryManifestV1 struct {
	SchemaVersion int                           `json:"schemaVersion"`
	Name          string                        `json:"name"`
	Tag           string                        `json:"tag"`
	Arch          string                        `json:"architecture"`
	FsLayers      []RegistryManifestV1Layer     `json:"fsLayers"`
	History       []RegistryManifestV1History   `json:"history"`
	Signatures    []RegistryManifestV1Signature `json:"signatures"`
}

// RegistryManifestV2 is v2 manifest spec
type RegistryManifestV2 struct {
	SchemaVersion int                      `json:"schemaVersion"`
	MediaType     string                   `json:"mediaType"`
	Config        *RegistryManifestV2Item  `json:"config"`
	Layers        []RegistryManifestV2Item `json:"layers"`
}

// RegistryManifestV2Item is manifest item
type RegistryManifestV2Item struct {
	MediaType string `json:"mediaType"`
	Size      int    `json:"size"`
	Digest    string `json:"digest"`
}
