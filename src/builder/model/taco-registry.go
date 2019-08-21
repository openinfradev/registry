package model

// BuildLogRow is building log 1 row
type BuildLogRow struct {
	BuildID string `json:"buildId"`
	Seq     int    `json:"seq"`
	Message string `json:"message"`
	Type    string `json:"type"`
	Valid   bool
}
