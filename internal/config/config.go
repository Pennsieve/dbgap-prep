package config

import (
	"fmt"
	"strings"
)

type ConsentGroup int

const (
	GRU ConsentGroup = iota
	HMB
	OTHER
)

func (g ConsentGroup) String() string {
	switch g {
	case GRU:
		return "GRU"
	case HMB:
		return "HMB"
	case OTHER:
		return "Other"
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
	case strings.ToLower(OTHER.String()):
		return OTHER, nil
	default:
		return 0, fmt.Errorf("unknown consent group '%s'", s)
	}
}

type AnalyteType int

const (
	DNA AnalyteType = iota
	RNA
	DNARNA
)

func (at AnalyteType) String() string {
	switch at {
	case DNA:
		return "DNA"
	case RNA:
		return "RNA"
	case DNARNA:
		return "DNA/RNA"
	default:
		return "UNKNOWN"
	}
}

func AnalyteTypeFromString(s string) (AnalyteType, error) {
	switch strings.ToLower(s) {
	case strings.ToLower(DNA.String()):
		return DNA, nil
	case strings.ToLower(RNA.String()):
		return RNA, nil
	case strings.ToLower(DNARNA.String()):
		return DNARNA, nil
	default:
		return 0, fmt.Errorf("unknown analyte type '%s'", s)
	}
}

type Config struct {
	IntegrationID      string
	WorkflowInstanceID string
	InputDirectory     string
	OutputDirectory    string
	ConsentGroup       ConsentGroup
	AnalyteType        AnalyteType
	IsTumor            bool
}
