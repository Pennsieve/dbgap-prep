package app

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/pennsieve/dbgap-prep/internal/config"
	"github.com/pennsieve/dbgap-prep/internal/enums/analytetype"
	"github.com/pennsieve/dbgap-prep/internal/enums/consentgroup"
	"github.com/pennsieve/dbgap-prep/internal/enums/istumor"
)

const ExecutionRunIDKey = "EXECUTION_RUN_ID"

// Older names for EXECUTION_RUN_ID, still set by the platform for processors that read them.
const IntegrationIDKey = "INTEGRATION_ID"
const WorkflowInstanceIDKey = "WORKFLOW_INSTANCE_ID"
const InputDirectoryKey = "INPUT_DIR"
const OutputDirectoryKey = "OUTPUT_DIR"
const ConsentGroupKey = "CONSENT_GROUP"
const AnalyteTypeKey = "ANALYTE_TYPE"
const IsTumorKey = "IS_TUMOR"
const PHSAccessionKey = "PHS_ACCESSION"

func ConfigFromEnv() (*config.Config, error) {
	var cfg config.Config
	var err error
	// Only used in logs, so not required.
	cfg.ExecutionRunID = LookupExecutionRunIDEnvVar()
	cfg.InputDirectory, err = LookupRequiredEnvVar(InputDirectoryKey)
	if err != nil {
		return nil, err
	}
	cfg.OutputDirectory, err = LookupRequiredEnvVar(OutputDirectoryKey)
	if err != nil {
		return nil, err
	}
	cfg.ConsentGroup, err = LookupConsentGroupEnvVar()
	if err != nil {
		return nil, err
	}
	cfg.AnalyteType, err = LookupAnalyteTypeEnvVar()
	if err != nil {
		return nil, err
	}
	cfg.IsTumor, err = LookupIsTumorEnvVar(istumor.NO)
	if err != nil {
		return nil, err
	}
	cfg.PHSAccession = os.Getenv(PHSAccessionKey)
	return &cfg, nil
}

// LookupExecutionRunIDEnvVar returns the run's ID from EXECUTION_RUN_ID, falling back to the
// older INTEGRATION_ID and WORKFLOW_INSTANCE_ID, or "" if none is set.
func LookupExecutionRunIDEnvVar() string {
	for _, key := range []string{ExecutionRunIDKey, IntegrationIDKey, WorkflowInstanceIDKey} {
		if value := os.Getenv(key); len(value) > 0 {
			return value
		}
	}
	return ""
}

func LookupRequiredEnvVar(key string) (string, error) {
	value := os.Getenv(key)
	if len(value) == 0 {
		return "", fmt.Errorf("no %s set", key)
	}
	return value, nil
}

func LookupConsentGroupEnvVar() (consentgroup.Group, error) {
	strValue, err := LookupRequiredEnvVar(ConsentGroupKey)
	if err != nil {
		return 0, err
	}
	return consentgroup.FromString(strValue)
}

func LookupAnalyteTypeEnvVar() (analytetype.Type, error) {
	strValue, err := LookupRequiredEnvVar(AnalyteTypeKey)
	if err != nil {
		return 0, err
	}
	return analytetype.FromString(strValue)
}

func LookupIsTumorEnvVar(defaultValue istumor.Value) (istumor.Value, error) {
	value := strings.TrimSpace(os.Getenv(IsTumorKey))
	if len(value) == 0 {
		logger.Info("env var not set; using default",
			slog.String("key", IsTumorKey),
			slog.String("default", defaultValue.String()),
		)
		return defaultValue, nil
	}
	return istumor.FromString(value)
}
