package url

// PathRegistryCatalog is registry catalog path
const PathRegistryCatalog string = "/v2/_catalog"

// PathRegistryTagList is repository tag list path
const PathRegistryTagList string = "/v2/%s/tags/list"

// PathRegistryManifest is repository manifests path
const PathRegistryManifest string = "/v2/%s/manifests/%s"

// PathRegistryBlobs is repository layer blob path
const PathRegistryBlobs string = "/v2/%s/blobs/%s"

// GitRepositoryPrivateURL is git repository url
const GitRepositoryPrivateURL string = "%s://%s:%s@%s"

// GitRepositoryPublicURL is git repository url
const GitRepositoryPublicURL string = "%s://%s"

// SecurityScanLayer is Clair
const SecurityScanLayer string = "http://%s/v1/layers/%s"

// SecurityScan is Clair
const SecurityScan string = "http://%s/v1/layers"

// SecurityScanLayerParam is Clair extra param
const SecurityScanLayerParam string = "features&vulnerabilities"
