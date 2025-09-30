package main

import (
	"flag"
	app "github.com/pennsieve/dbgap-prep/internal"
	"github.com/pennsieve/dbgap-prep/internal/logging"
	"log/slog"
	"os"
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

	dbgap := app.NewApp("NA", inputDirectory, outputDirectory)

	logger.Info("created local dbgap-prep application",
		slog.String("integrationID", dbgap.IntegrationID),
		slog.String("inputDirectory", dbgap.InputDirectory),
		slog.String("outputDirectory", dbgap.OutputDirectory),
	)

	if err := dbgap.Run(); err != nil {
		logger.Error("error running local dbgap-prep application", slog.Any("error", err))
		os.Exit(1)
	}
}
