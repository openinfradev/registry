package scope

import "fmt"

// ScopeTemplate is scope pattern template
const ScopeTemplate string = "%s:%s:%s"

// TypeRegistry is scope registry type
const TypeRegistry string = "registry"

// TypeRepository is scope repository type
const TypeRepository string = "repository"

// ResourceCatalog is res catalog
const ResourceCatalog string = "catalog"

// ResourceRepository is res repository
const ResourceRepository string = "repository"

// ResourceManifest is res manifest
const ResourceManifest string = "manifest"

// ResourceBlob is res blob
const ResourceBlob string = "blob"

// ActionWildCard is action wildcard
const ActionWildCard string = "*"

// Scope is scope struct
type Scope struct {
	Type     string
	Resource string
	Action   string
}

// String returns scope pattern string
func (s *Scope) String() string {
	return fmt.Sprintf(ScopeTemplate, s.Type, s.Resource, s.Action)
}
