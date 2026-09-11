package app

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/pennsieve/dbgap-prep/internal/config"
	"github.com/pennsieve/dbgap-prep/internal/enums/analytetype"
	"github.com/pennsieve/dbgap-prep/internal/enums/consentgroup"
	"github.com/pennsieve/dbgap-prep/internal/enums/istumor"
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
	cfg.IsTumor, err = LookupIsTumorEnvVar(IsTumorKey, istumor.NO)
	return &cfg, nil
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

func LookupIsTumorEnvVar(key string, defaultValue istumor.Value) (istumor.Value, error) {
	value := os.Getenv(key)
	if len(value) == 0 {
		logger.Info("env var not set; using default",
			slog.String("key", key),
			slog.String("default", defaultValue.String()),
		)
		return defaultValue, nil
	}
	return istumor.FromString(value)
}
