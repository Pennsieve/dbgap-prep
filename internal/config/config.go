package config

import (
	"github.com/pennsieve/dbgap-prep/internal/enums/analytetype"
	"github.com/pennsieve/dbgap-prep/internal/enums/consentgroup"
	"github.com/pennsieve/dbgap-prep/internal/enums/istumor"
)

type Config struct {
	IntegrationID      string
	WorkflowInstanceID string
	InputDirectory     string
	OutputDirectory    string
	ConsentGroup       consentgroup.Group
	AnalyteType        analytetype.Type
	IsTumor            istumor.Value
}
