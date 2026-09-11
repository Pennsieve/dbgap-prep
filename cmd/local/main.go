package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	app "github.com/pennsieve/dbgap-prep/internal"
	"github.com/pennsieve/dbgap-prep/internal/config"
	"github.com/pennsieve/dbgap-prep/internal/logging"
)

var logger = logging.PackageLogger("main")

var inputDirectory string
var outputDirectory string
var consentGroup string
var analyteType string
var isTumor bool

func init() {
	inputUsage := "input directory containing dataset_description.xlsx, subjects.xlsx, and samples.xlsx"
	flag.StringVar(&inputDirectory, "input-directory", "", inputUsage)
	flag.StringVar(&inputDirectory, "i", "", inputUsage+" (shorthand)")

	outputUsage := "output director where dbGaP files will be written"
	flag.StringVar(&outputDirectory, "output-directory", "", outputUsage)
	flag.StringVar(&outputDirectory, "o", "", outputUsage+" (shorthand)")

	consentGroupUsage := fmt.Sprintf("consent group name; either %s or %s", config.GRU, config.HMB)
	flag.StringVar(&consentGroup, "consent-group", "", consentGroupUsage)
	flag.StringVar(&consentGroup, "c", "", consentGroupUsage+" (shorthand)")

	analyteTypeUsage := fmt.Sprintf("analyte type; one of %s, %s, %s", config.DNA, config.RNA, config.DNARNA)
	flag.StringVar(&analyteType, "analyte-type", "", analyteTypeUsage)
	flag.StringVar(&analyteType, "a", "", analyteTypeUsage+" (shorthand)")

	isTumorUsage := "true if samples are tumors; specified as -t (or -t=false)"
	flag.BoolVar(&isTumor, "tumor", false, isTumorUsage)
	flag.BoolVar(&isTumor, "t", false, isTumorUsage+" (shorthand)")
}

func main() {
	flag.Parse()

	if len(inputDirectory) == 0 {
		logger.Error("missing input directory")
		flag.Usage()
		os.Exit(1)
	}

	if len(consentGroup) == 0 {
		logger.Error("missing consent group")
		flag.Usage()
		os.Exit(1)
	}

	if len(analyteType) == 0 {
		logger.Error("missing analyte type")
		flag.Usage()
		os.Exit(1)
	}

	cg, err := config.ConsentGroupFromString(consentGroup)
	if err != nil {
		logger.Error(err.Error())
		flag.Usage()
		os.Exit(1)
	}

	at, err := config.AnalyteTypeFromString(analyteType)
	if err != nil {
		logger.Error(err.Error())
		flag.Usage()
		os.Exit(1)
	}

	cfg := &config.Config{
		IntegrationID:      "NA",
		WorkflowInstanceID: "NA",
		InputDirectory:     inputDirectory,
		OutputDirectory:    outputDirectory,
		ConsentGroup:       cg,
		AnalyteType:        at,
		IsTumor:            isTumor,
	}

	dbgap := app.NewApp(cfg)

	logger.Info("created local dbgap-prep application",
		slog.String("integrationID", dbgap.Config.IntegrationID),
		slog.String("inputDirectory", dbgap.Config.InputDirectory),
		slog.String("outputDirectory", dbgap.Config.OutputDirectory),
		slog.String("consentGroup", dbgap.Config.ConsentGroup.String()),
		slog.String("analyteType", dbgap.Config.AnalyteType.String()),
		slog.Bool("isTumor", dbgap.Config.IsTumor),
	)

	if err := dbgap.Run(); err != nil {
		logger.Error("error running local dbgap-prep application", slog.Any("error", err))
		os.Exit(1)
	}
}
