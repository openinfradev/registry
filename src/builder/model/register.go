package model

// Builder is builder host information
type Builder struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

// BuilderList is builder host list
type BuilderList struct {
	Builders []Builder `json:"builders"`
}
