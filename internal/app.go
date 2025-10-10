package app

import (
	"fmt"
	"github.com/pennsieve/dbgap-prep/internal/dbgap/sampleattributes"
	"github.com/pennsieve/dbgap-prep/internal/dbgap/subjectconsent"
	scds "github.com/pennsieve/dbgap-prep/internal/dbgap/subjectconsent/ds"
	"github.com/pennsieve/dbgap-prep/internal/dbgap/subjectsample"
	"github.com/pennsieve/dbgap-prep/internal/logging"
	"github.com/pennsieve/dbgap-prep/internal/samples"
	"github.com/pennsieve/dbgap-prep/internal/subjects"
	"github.com/pennsieve/dbgap-prep/internal/utils"
	"github.com/xuri/excelize/v2"
	"log/slog"
	"path/filepath"
)

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

	samplesPath := filepath.Join(a.InputDirectory, samples.FileName)
	samplesLogger := logger.With(slog.String("file", samplesPath))
	samplesFile, err := openExcelInput(samplesPath)
	if err != nil {
		return err
	}
	defer utils.CloseExcelFile(samplesFile, samplesLogger)

	samplesLogger.Info("reading samples file")

	samplesHeader, samps, err := samples.FromFile(samplesFile)
	if err != nil {
		return err
	}

	consentedSubjects, consentedSamples, consentedSamplesInSubjectOrder := scds.GetConsented(subjectsConsents, samps)
	logger.Info("filtered by consent",
		slog.Int("totalSubjects", len(subs)),
		slog.Int("consentedSubjects", len(consentedSubjects)),
		slog.Int("totalSamples", len(samps)),
		slog.Int("consentedSamples", len(consentedSamples)),
	)

	//TODO write subject phenotype files

	if len(consentedSamples) == 0 {
		samplesLogger.Info("no consented samples found")
		return nil
	}

	if err := subjectsample.WriteFiles(a.OutputDirectory, consentedSamplesInSubjectOrder); err != nil {
		return err
	}

	if err := sampleattributes.WriteFiles(a.OutputDirectory, samplesHeader, consentedSamples); err != nil {
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
