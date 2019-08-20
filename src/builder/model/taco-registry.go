package model

import (
	"regexp"
)

// BuildLogRow is building log 1 row
type BuildLogRow struct {
	BuildID string `json:"buildId"`
	Seq     int    `json:"seq"`
	Message string `json:"message"`
	Type    string `json:"type"`
}

// Parse is raw log parsing
func (b *BuildLogRow) Parse(buildID string, seq int, raw string) {
	m, _ := regexp.MatchString("(Step)\\s+[0-9]+(/)[0-9]+\\s+(:)", raw)
	if m {
		b.Type = "command"
	}
	b.BuildID = buildID
	b.Seq = seq
	b.Message = raw
}
