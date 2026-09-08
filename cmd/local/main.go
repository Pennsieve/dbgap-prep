package main

import (
	"flag"
	"log/slog"
	"os"

	app "github.com/pennsieve/dbgap-prep/internal"
	"github.com/pennsieve/dbgap-prep/internal/logging"
)

var logger = logging.PackageLogger("main")

var inputDirectory string
var outputDirectory string

func init() {
	inputUsage := "input directory containing subjects.xlsx and samples.xlsx"
	flag.StringVar(&inputDirectory, "input-directory", "", inputUsage)
	flag.StringVar(&inputDirectory, "i", "", inputUsage+" (shorthand)")

	outputUsage := "output director where dbGaP files will be written"
	flag.StringVar(&outputDirectory, "output-directory", "", outputUsage)
	flag.StringVar(&outputDirectory, "o", "", outputUsage+" (shorthand)")
}
func main() {
	flag.Parse()

	if len(inputDirectory) == 0 {
		logger.Error("missing input directory")
		flag.Usage()
		os.Exit(1)
	}

	config := &app.Config{
		IntegrationID:      "NA",
		WorkflowInstanceID: "NA",
		InputDirectory:     inputDirectory,
		OutputDirectory:    outputDirectory,
	}

	dbgap := app.NewApp(config)

	logger.Info("created local dbgap-prep application",
		slog.String("integrationID", dbgap.Config.IntegrationID),
		slog.String("inputDirectory", dbgap.Config.InputDirectory),
		slog.String("outputDirectory", dbgap.Config.OutputDirectory),
	)

	if err := dbgap.Run(); err != nil {
		logger.Error("error running local dbgap-prep application", slog.Any("error", err))
		os.Exit(1)
	}
}
