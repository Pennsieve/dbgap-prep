package lambda

import (
	"context"
	"fmt"
	"log/slog"

	app "github.com/pennsieve/dbgap-prep/internal"
	"github.com/pennsieve/dbgap-prep/internal/logging"
	"github.com/pennsieve/dbgap-prep/internal/utils"
)

var logger = logging.PackageLogger("lambda")

type Event struct {
	IntegrationID      string         `json:"integrationId"`
	WorkflowInstanceID string         `json:"workflowInstanceId"`
	InputDirectory     string         `json:"inputDir"`
	OutputDirectory    string         `json:"outputDir"`
	ConsentGroup       string         `json:"CONSENT_GROUP"`
	AnalyteType        string         `json:"ANALYTE_TYPE"`
	IsTumor            utils.FlexBool `json:"IS_TUMOR"`
}

func Handler(_ context.Context, event Event) error {
	logger.Info("lambda handler invoked",
		slog.String("integrationID", event.IntegrationID),
		slog.String("workflowInstanceID", event.WorkflowInstanceID),
		slog.String("inputDirectory", event.InputDirectory),
		slog.String("outputDirectory", event.OutputDirectory),
		slog.String("consentGroup", event.ConsentGroup),
		slog.String("analyteType", event.AnalyteType),
		slog.Bool("isTumor", bool(event.IsTumor)),
	)

	m := app.NewApp(configFromEvent(event))

	if err := m.Run(); err != nil {
		return fmt.Errorf("error running application: %w", err)
	}

	return nil
}

func configFromEvent(event Event) *app.Config {
	return &app.Config{
		IntegrationID:      event.IntegrationID,
		WorkflowInstanceID: event.WorkflowInstanceID,
		InputDirectory:     event.InputDirectory,
		OutputDirectory:    event.OutputDirectory,
		ConsentGroup:       event.ConsentGroup,
		AnalyteType:        event.AnalyteType,
		IsTumor:            bool(event.IsTumor),
	}
}
