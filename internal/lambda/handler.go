package lambda

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	app "github.com/pennsieve/dbgap-prep/internal"
	"github.com/pennsieve/dbgap-prep/internal/config"
	"github.com/pennsieve/dbgap-prep/internal/enums/analytetype"
	"github.com/pennsieve/dbgap-prep/internal/enums/consentgroup"
	"github.com/pennsieve/dbgap-prep/internal/enums/istumor"
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
	PHSAccession       string `json:"PHS_ACCESSION"`
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
		slog.String("phsAccession", event.PHSAccession),
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
	consentGroup, err := consentgroup.FromString(event.ConsentGroup)
	if err != nil {
		return nil, err
	}
	analyteType, err := analytetype.FromString(event.AnalyteType)
	if err != nil {
		return nil, err
	}
	isTumor := istumor.NO
	if isTumorString := strings.TrimSpace(event.IsTumor); len(isTumorString) > 0 {
		isTumor, err = istumor.FromString(isTumorString)
		if err != nil {
			return nil, err
		}
	}
	return &config.Config{
		IntegrationID:      event.IntegrationID,
		WorkflowInstanceID: event.WorkflowInstanceID,
		InputDirectory:     event.InputDirectory,
		OutputDirectory:    event.OutputDirectory,
		ConsentGroup:       consentGroup,
		AnalyteType:        analyteType,
		IsTumor:            isTumor,
		PHSAccession:       event.PHSAccession,
	}, nil
}
