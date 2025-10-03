package app

import (
	"fmt"
	"github.com/pennsieve/dbgap-prep/internal/dbgap/subjectconsent"
	"github.com/pennsieve/dbgap-prep/internal/dbgap/subjectsample"
	"github.com/pennsieve/dbgap-prep/internal/logging"
	"github.com/pennsieve/dbgap-prep/internal/samples"
	"github.com/pennsieve/dbgap-prep/internal/subjects"
	"github.com/pennsieve/dbgap-prep/internal/utils"
	"github.com/xuri/excelize/v2"
	"log/slog"
	"path/filepath"
)

const samplesFileName = "samples.xlsx"

var logger = logging.PackageLogger("app")

type App struct {
	IntegrationID   string
	InputDirectory  string
	OutputDirectory string
}

func NewApp(integrationID string, inputDirectory string, outputDirectory string) *App {
	return &App{
		IntegrationID:   integrationID,
		InputDirectory:  inputDirectory,
		OutputDirectory: outputDirectory,
	}
}

func (a *App) Run() error {
	subjectsPath := filepath.Join(a.InputDirectory, subjects.FileName)
	subjectsLogger := logger.With(slog.String("file", subjectsPath))

	subjectsFile, err := openExcelInput(subjectsPath)
	if err != nil {
		return err
	}
	defer utils.CloseExcelFile(subjectsFile, subjectsLogger)

	subjectsLogger.Info("reading subjects file")

	subs, err := subjects.FromFile(subjectsFile)
	if err != nil {
		return err
	}
	if len(subs) == 0 {
		subjectsLogger.Info("no subjects found; no dbGaP files created")
		// Nothing to do. Or should we create empty dbGaP files?
		return nil
	}

	subjectsConsents, err := subjectconsent.WriteFiles(a.OutputDirectory, subs)
	if err != nil {
		return err
	}

	samplesPath := filepath.Join(a.InputDirectory, samplesFileName)
	samplesLogger := logger.With(slog.String("file", samplesPath))
	samplesFile, err := openExcelInput(samplesPath)
	if err != nil {
		return err
	}
	defer utils.CloseExcelFile(samplesFile, samplesLogger)

	samplesLogger.Info("reading samples file")

	samps, err := samples.FromFile(samplesFile)
	if err != nil {
		return err
	}

	if len(samps) == 0 {
		samplesLogger.Info("no samples found")
		return nil
	}

	if err := subjectsample.WriteFiles(a.OutputDirectory, subjectsConsents, samps); err != nil {
		return err
	}

	return nil
}

func openExcelInput(filePath string) (*excelize.File, error) {
	inputFile, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error opening input file %s: %w", filePath, err)
	}
	return inputFile, nil
}
