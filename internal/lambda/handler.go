package lambda

import (
	"context"
	"fmt"
	app "github.com/pennsieve/dbgap-prep/internal"
	"github.com/pennsieve/dbgap-prep/internal/logging"
	"log/slog"
)

var logger = logging.PackageLogger("lambda")

type Event struct {
	IntegrationID      string `json:"integrationId"`
	WorkflowInstanceID string `json:"workflowInstanceId"`
	InputDirectory     string `json:"inputDir"`
	OutputDirectory    string `json:"outputDir"`
}

func Handler(_ context.Context, event Event) error {
	logger.Info("lambda handler invoked",
		slog.String("integrationID", event.IntegrationID),
		slog.String("workflowInstanceID", event.WorkflowInstanceID),
		slog.String("inputDirectory", event.InputDirectory),
		slog.String("outputDirectory", event.OutputDirectory),
	)

	m := app.NewApp(event.IntegrationID, event.WorkflowInstanceID, event.InputDirectory, event.OutputDirectory)

	if err := m.Run(); err != nil {
		return fmt.Errorf("error running application: %w", err)
	}

	return nil
}
