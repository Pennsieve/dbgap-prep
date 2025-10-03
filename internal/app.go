package app

import (
	"fmt"
	"github.com/pennsieve/dbgap-prep/internal/dbgap/dd"
	scdd "github.com/pennsieve/dbgap-prep/internal/dbgap/subjectconsent/dd"
	scds "github.com/pennsieve/dbgap-prep/internal/dbgap/subjectconsent/ds"
	"github.com/pennsieve/dbgap-prep/internal/logging"
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
	subjectConsentDDPath := filepath.Join(a.OutputDirectory, scdd.Spec.FileName)

	if err := dd.Write(subjectConsentDDPath, scdd.Spec); err != nil {
		return err
	}

	logger.Info("wrote subject consent DD file", slog.String("file", subjectConsentDDPath))

	subjectConsentDSPath := filepath.Join(a.OutputDirectory, scds.Spec.FileName)
	if err := scds.Write(subjectConsentDSPath, subs); err != nil {
		return err
	}

	logger.Info("wrote subject consent DS file", slog.String("file", subjectConsentDSPath))

	samplesPath := filepath.Join(a.InputDirectory, samplesFileName)
	samplesLogger := logger.With(slog.String("file", samplesPath))
	samples, err := openExcelInput(samplesPath)
	if err != nil {
		return err
	}
	defer utils.CloseExcelFile(samples, samplesLogger)

	samplesLogger.Info("reading samples file")

	return nil
}

func openExcelInput(filePath string) (*excelize.File, error) {
	inputFile, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error opening input file %s: %w", filePath, err)
	}
	return inputFile, nil
}
