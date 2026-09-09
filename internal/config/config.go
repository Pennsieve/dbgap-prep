package config

import (
	"fmt"
	"strings"
)

type ConsentGroup int

const (
	GRU ConsentGroup = iota
	HMB
)

func (g ConsentGroup) String() string {
	switch g {
	case GRU:
		return "GRU"
	case HMB:
		return "HMB"
	default:
		return "UNKNOWN"
	}
}

func ConsentGroupFromString(s string) (ConsentGroup, error) {
	switch strings.ToLower(s) {
	case strings.ToLower(GRU.String()):
		return GRU, nil
	case strings.ToLower(HMB.String()):
		return HMB, nil
	default:
		return 0, fmt.Errorf("unknown consent group '%s'", s)
	}
}

type Config struct {
	IntegrationID      string
	WorkflowInstanceID string
	InputDirectory     string
	OutputDirectory    string
	ConsentGroup       ConsentGroup
	AnalyteType        string
	IsTumor            bool
}
