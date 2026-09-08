package main

import (
	"log/slog"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	app "github.com/pennsieve/dbgap-prep/internal"
	lambdahandler "github.com/pennsieve/dbgap-prep/internal/lambda"
	"github.com/pennsieve/dbgap-prep/internal/logging"
)

var logger = logging.PackageLogger("main")

func main() {

	if _, isLambda := os.LookupEnv("AWS_LAMBDA_RUNTIME_API"); isLambda {
		lambda.Start(lambdahandler.Handler)
		return
	}

	config, err := app.ConfigFromEnv()
	if err != nil {
		logger.Error("error loading config from environment", slog.Any("error", err))
	}

	logger.Info("loaded config from environment",
		slog.String("integrationID", config.IntegrationID),
		slog.String("workflowInstanceID", config.WorkflowInstanceID),
		slog.String("inputDirectory", config.InputDirectory),
		slog.String("outputDirectory", config.OutputDirectory),
		slog.String("consentGroup", config.ConsentGroup),
		slog.String("analyteType", config.AnalyteType),
		slog.Bool("isTumor", config.IsTumor),
	)

	m, err := app.FromEnv()
	if err != nil {
		logger.Error("error creating application", slog.Any("error", err))
		os.Exit(1)
	}

	logger.Info("created dbgap-prep application in ECS mode",
		slog.String("integrationID", m.IntegrationID),
		slog.String("workflowInstanceID", m.WorkflowInstanceID),
		slog.String("inputDirectory", m.InputDirectory),
		slog.String("outputDirectory", m.OutputDirectory),
	)

	if err := m.Run(); err != nil {
		logger.Error("error running application", slog.Any("error", err))
		os.Exit(1)
	}
}
