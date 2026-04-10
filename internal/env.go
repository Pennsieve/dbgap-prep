package app

import (
	"fmt"
	"os"
)

const IntegrationIDKey = "INTEGRATION_ID"
const WorkflowInstanceID = "WORKFLOW_INSTANCE_ID"
const InputDirectoryKey = "INPUT_DIR"
const OutputDirectoryKey = "OUTPUT_DIR"

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

func LookupRequiredEnvVar(key string) (string, error) {
	value := os.Getenv(key)
	if len(value) == 0 {
		return "", fmt.Errorf("no %s set", key)
	}
	return value, nil
}
