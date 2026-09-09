package app

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"

	"github.com/pennsieve/dbgap-prep/internal/config"
)

const IntegrationIDKey = "INTEGRATION_ID"
const WorkflowInstanceID = "WORKFLOW_INSTANCE_ID"
const InputDirectoryKey = "INPUT_DIR"
const OutputDirectoryKey = "OUTPUT_DIR"
const ConsentGroupKey = "CONSENT_GROUP"
const AnalyteTypeKey = "ANALYTE_TYPE"
const IsTumorKey = "IS_TUMOR"

func ConfigFromEnv() (*config.Config, error) {
	var cfg config.Config
	var err error
	cfg.IntegrationID, err = LookupRequiredEnvVar(IntegrationIDKey)
	if err != nil {
		return nil, err
	}
	// Not clear if this will be present, so not required.
	cfg.WorkflowInstanceID = os.Getenv(WorkflowInstanceID)
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
	cfg.IsTumor, err = LookupBoolEnvVar(IsTumorKey, false)
	return &cfg, nil
}

func LookupRequiredEnvVar(key string) (string, error) {
	value := os.Getenv(key)
	if len(value) == 0 {
		return "", fmt.Errorf("no %s set", key)
	}
	return value, nil
}

func LookupConsentGroupEnvVar() (config.ConsentGroup, error) {
	strValue, err := LookupRequiredEnvVar(ConsentGroupKey)
	if err != nil {
		return 0, err
	}
	return config.ConsentGroupFromString(strValue)
}

func LookupAnalyteTypeEnvVar() (config.AnalyteType, error) {
	strValue, err := LookupRequiredEnvVar(AnalyteTypeKey)
	if err != nil {
		return 0, err
	}
	return config.AnalyteTypeFromString(strValue)
}

func LookupBoolEnvVar(key string, defaultValue bool) (bool, error) {
	value := os.Getenv(key)
	if len(value) == 0 {
		logger.Info("env var not set; using default",
			slog.String("key", key),
			slog.Bool("default", defaultValue),
		)
		return defaultValue, nil
	}
	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		// Thinking that if there is an unparsable value,
		// we should treat as an error rather than silently use the default.
		return false, err
	}
	return boolValue, nil
}
