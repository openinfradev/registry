package taco

import (
	"github.com/openinfradev/registry-builder/model"
	"regexp"
)

// ParseLog is raw log parsing
func ParseLog(buildID string, seq int, raw string) *model.BuildLogRow {
	b := &model.BuildLogRow{}

	v, _ := regexp.MatchString("^(Sending build context to Docker daemon)", raw)
	b.Valid = !v

	m, _ := regexp.MatchString("^(Step)\\s+[0-9]+(/)[0-9]+\\s+(:)", raw)
	if m {
		b.Type = "command"
	}
	b.BuildID = buildID
	b.Seq = seq
	b.Message = raw

	return b
}

// MakePhaseLog returns phase log
func MakePhaseLog(buildID string, seq int, raw string) *model.BuildLogRow {
	return &model.BuildLogRow{
		Type:    "phase",
		BuildID: buildID,
		Seq:     seq,
		Message: raw,
		Valid:   true,
	}
}
