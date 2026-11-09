package main

import (
	"log/slog"
	"os"

	app "github.com/pennsieve/dbgap-prep/internal"
	"github.com/pennsieve/dbgap-prep/internal/logging"
)

var logger = logging.PackageLogger("main")

// main runs the app once. It runs the same way as an ECS task and as a Lambda step: on Lambda,
// the platform's Lambda step runtime runs this entrypoint with the same environment variables.
func main() {
	config, err := app.ConfigFromEnv()
	if err != nil {
		logger.Error("error loading config from environment", slog.Any("error", err))
		os.Exit(1)
	}

	logger.Info("loaded config from environment",
		slog.String("executionRunID", config.ExecutionRunID),
		slog.String("inputDirectory", config.InputDirectory),
		slog.String("outputDirectory", config.OutputDirectory),
		slog.String("consentGroup", config.ConsentGroup.String()),
		slog.String("analyteType", config.AnalyteType.String()),
		slog.String("isTumor", config.IsTumor.String()),
	)

	m := app.NewApp(config)

	logger.Info("created dbgap-prep application")

	if err := m.Run(); err != nil {
		logger.Error("error running application", slog.Any("error", err))
		os.Exit(1)
	}
}
