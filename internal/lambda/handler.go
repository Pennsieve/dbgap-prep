package lambda

import (
	"context"
	"fmt"
	"log/slog"

	app "github.com/pennsieve/dbgap-prep/internal"
	"github.com/pennsieve/dbgap-prep/internal/config"
	"github.com/pennsieve/dbgap-prep/internal/logging"
)

var logger = logging.PackageLogger("lambda")

type Event struct {
	IntegrationID      string `json:"integrationId"`
	WorkflowInstanceID string `json:"workflowInstanceId"`
	InputDirectory     string `json:"inputDir"`
	OutputDirectory    string `json:"outputDir"`
	ConsentGroup       string `json:"CONSENT_GROUP"`
	AnalyteType        string `json:"ANALYTE_TYPE"`
	IsTumor            string `json:"IS_TUMOR"`
}

func Handler(_ context.Context, event Event) error {
	logger.Info("lambda handler invoked",
		slog.String("integrationID", event.IntegrationID),
		slog.String("workflowInstanceID", event.WorkflowInstanceID),
		slog.String("inputDirectory", event.InputDirectory),
		slog.String("outputDirectory", event.OutputDirectory),
		slog.String("consentGroup", event.ConsentGroup),
		slog.String("analyteType", event.AnalyteType),
		slog.String("isTumor", event.IsTumor),
	)

	cfg, err := configFromEvent(event)
	if err != nil {
		return fmt.Errorf("error converting event to config: %w", err)
	}

	m := app.NewApp(cfg)

	if err := m.Run(); err != nil {
		return fmt.Errorf("error running application: %w", err)
	}

	return nil
}

func configFromEvent(event Event) (*config.Config, error) {
	consentGroup, err := config.ConsentGroupFromString(event.ConsentGroup)
	if err != nil {
		return nil, err
	}
	analyteType, err := config.AnalyteTypeFromString(event.AnalyteType)
	if err != nil {
		return nil, err
	}
	isTumor, err := config.IsTumorFromString(event.IsTumor)
	if err != nil {
		return nil, err
	}
	return &config.Config{
		IntegrationID:      event.IntegrationID,
		WorkflowInstanceID: event.WorkflowInstanceID,
		InputDirectory:     event.InputDirectory,
		OutputDirectory:    event.OutputDirectory,
		ConsentGroup:       consentGroup,
		AnalyteType:        analyteType,
		IsTumor:            isTumor,
	}, nil
}
