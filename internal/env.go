package app

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"
)

const IntegrationIDKey = "INTEGRATION_ID"
const WorkflowInstanceID = "WORKFLOW_INSTANCE_ID"
const InputDirectoryKey = "INPUT_DIR"
const OutputDirectoryKey = "OUTPUT_DIR"
const ConsentGroupKey = "CONSENT_GROUP"
const AnalyteTypeKey = "ANALYTE_TYPE"
const IsTumorKey = "IS_TUMOR"

func FromEnv() (*App, error) {
	integrationID, err := LookupRequiredEnvVar(IntegrationIDKey)
	if err != nil {
		return nil, err
	}
	// Not clear if this will be present, so not required.
	workflowInstanceID := os.Getenv(WorkflowInstanceID)
	inputDirectory, err := LookupRequiredEnvVar(InputDirectoryKey)
	if err != nil {
		return nil, err
	}
	outputDirectory, err := LookupRequiredEnvVar(OutputDirectoryKey)
	if err != nil {
		return nil, err
	}
	return NewApp(integrationID,
		workflowInstanceID,
		inputDirectory,
		outputDirectory,
	), nil
}

func ConfigFromEnv() (*Config, error) {
	var config Config
	var err error
	config.IntegrationID, err = LookupRequiredEnvVar(IntegrationIDKey)
	if err != nil {
		return nil, err
	}
	// Not clear if this will be present, so not required.
	config.WorkflowInstanceID = os.Getenv(WorkflowInstanceID)
	config.InputDirectory, err = LookupRequiredEnvVar(InputDirectoryKey)
	if err != nil {
		return nil, err
	}
	config.OutputDirectory, err = LookupRequiredEnvVar(OutputDirectoryKey)
	if err != nil {
		return nil, err
	}
	config.ConsentGroup, err = LookupRequiredEnvVar(ConsentGroupKey)
	if err != nil {
		return nil, err
	}
	config.AnalyteType, err = LookupRequiredEnvVar(AnalyteTypeKey)
	if err != nil {
		return nil, err
	}
	config.IsTumor, err = LookupBoolEnvVar(IsTumorKey, false)
	return &config, nil
}

func LookupRequiredEnvVar(key string) (string, error) {
	value := os.Getenv(key)
	if len(value) == 0 {
		return "", fmt.Errorf("no %s set", key)
	}
	return value, nil
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
