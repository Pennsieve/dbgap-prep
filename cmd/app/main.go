package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	app "github.com/pennsieve/dbgap-prep/internal"
	lambdahandler "github.com/pennsieve/dbgap-prep/internal/lambda"
	"github.com/pennsieve/dbgap-prep/internal/logging"
	"log/slog"
	"os"
)

var logger = logging.PackageLogger("main")

func main() {

	if _, isLambda := os.LookupEnv("AWS_LAMBDA_RUNTIME_API"); isLambda {
		lambda.Start(lambdahandler.Handler)
		return
	}

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
