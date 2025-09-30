package main

import (
	app "github.com/pennsieve/dbgap-prep/internal"
	"github.com/pennsieve/dbgap-prep/internal/logging"
	"log/slog"
	"os"
)

var logger = logging.PackageLogger("main")

func main() {
	m, err := app.FromEnv()
	if err != nil {
		logger.Error("error creating application", slog.Any("error", err))
		os.Exit(1)
	}

	logger.Info("created dbgap-prep application",
		slog.String("integrationID", m.IntegrationID),
		slog.String("inputDirectory", m.InputDirectory),
		slog.String("outputDirectory", m.OutputDirectory),
	)

	if err := m.Run(); err != nil {
		logger.Error("error running application", slog.Any("error", err))
		os.Exit(1)
	}
}
