package app

import (
	"fmt"
	"github.com/pennsieve/dbgap-prep/internal/logging"
	"github.com/xuri/excelize/v2"
	"log/slog"
	"path/filepath"
)

const subjectsFile = "subjects.xlsx"
const samplesFile = "samples.xlsx"

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
	subjectsPath := filepath.Join(a.InputDirectory, subjectsFile)
	subjects, err := openInput(subjectsPath)
	if err != nil {
		return err
	}
	defer closeInput(subjects)

	subjectSheets := subjects.GetSheetList()
	if len(subjectSheets) == 0 {
		return fmt.Errorf("no sheets found in %s", subjectsPath)
	}
	logger.Info("subject file sheets",
		slog.Int("sheetCount", len(subjectSheets)),
		slog.String("firstSheet", subjectSheets[0]))
	samplesPath := filepath.Join(a.InputDirectory, samplesFile)
	samples, err := openInput(samplesPath)
	if err != nil {
		return err
	}
	defer closeInput(samples)

	return nil
}

func openInput(filePath string) (*excelize.File, error) {
	inputFile, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error opening input file %s: %w", filePath, err)
	}
	return inputFile, nil
}

func closeInput(inputFile *excelize.File) {
	if err := inputFile.Close(); err != nil {
		logger.Warn("error closing input file",
			slog.String("path", inputFile.Path),
			slog.Any("error", err))
	}
}
