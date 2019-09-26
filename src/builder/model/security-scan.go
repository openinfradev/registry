package model

// SecurityScanLayer is security scanning layer result
type SecurityScanLayer struct {
	Status string                  `json:"status"`
	Data   *map[string]interface{} `json:"data"`
}

// SecurityScanParam is security scan parameters
type SecurityScanParam struct {
	Layer *SecurityScanLayerParam `json:"Layer"`
}

// SecurityScanLayerParam is security scan layer parameters
type SecurityScanLayerParam struct {
	Name       string                        `json:"Name"`
	ParentName string                        `json:"ParentName"`
	Path       string                        `json:"Path"`
	Format     string                        `json:"Format"`
	Headers    *SecurityScanLayerHeaderParam `json:"Headers"`
}

// SecurityScanLayerHeaderParam is authorization token
type SecurityScanLayerHeaderParam struct {
	Authorization string `json:"Authorization"`
}
