package app

import (
	"fmt"
	"github.com/pennsieve/dbgap-prep/internal/dbgap/sampleattributes"
	"github.com/pennsieve/dbgap-prep/internal/dbgap/subjectconsent"
	scds "github.com/pennsieve/dbgap-prep/internal/dbgap/subjectconsent/ds"
	"github.com/pennsieve/dbgap-prep/internal/dbgap/subjectphenotypes"
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

	subjectsHeader, subs, err := subjects.FromFile(subjectsFile)
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

	consentedSubjects, consentedSamples, consentedSamplesInSubjectOrder := scds.GetConsented(subjectsConsents, subs, samps)
	logger.Info("filtered by consent",
		slog.Int("totalSubjects", len(subs)),
		slog.Int("consentedSubjects", len(consentedSubjects)),
		slog.Int("totalSamples", len(samps)),
		slog.Int("consentedSamples", len(consentedSamples)),
	)

	// prune the subjects header of empty columns so our subject phenotype DD file does not contain
	// empty columns.
	subjectsHeader = pruneHeader(subjectsHeader, consentedSubjects, subjects.IDLabel, subjects.SexLabel)

	if err := subjectphenotypes.WriteFiles(a.OutputDirectory, subjectsHeader, consentedSubjects); err != nil {
		return err
	}

	if len(consentedSamples) == 0 {
		samplesLogger.Info("no consented samples found")
		return nil
	}

	if err := subjectsample.WriteFiles(a.OutputDirectory, consentedSamplesInSubjectOrder); err != nil {
		return err
	}

	// prune the samples header of empty columns so our samples attributes DD file does not contain
	// empty columns.
	samplesHeader = pruneHeader(samplesHeader, consentedSamples, samples.IDLabel, samples.SubjectIDLabel)
	if err := sampleattributes.WriteFiles(a.OutputDirectory, samplesHeader, consentedSamples); err != nil {
		return err
	}

	return nil
}

type HasValues interface {
	GetValue(key string) (string, bool)
}

// pruneHeader returns a string slice that is a subset of header, preserving the order. The returned
// slice will retain those header values that correspond to at least one non-empty string value in the given slice of HasValues.
// That is, it is a way of removing header labels that correspond to empty columns. If a header label should be retained even
// if its column is empty, pass it as a alwaysKeep value.
func pruneHeader[T HasValues](header []string, haveValues []T, alwaysKeep ...string) []string {
	keepers := make(map[string]bool, len(header))
	for _, keeper := range alwaysKeep {
		keepers[keeper] = true
	}

	for _, hv := range haveValues {
		for _, label := range header {
			if !keepers[label] {
				if value, ok := hv.GetValue(label); ok && len(value) > 0 {
					keepers[label] = true
				}
			}
		}
		if len(keepers) == len(header) {
			break // all headers accounted for
		}
	}
	pruned := make([]string, 0, len(keepers))
	for _, label := range header {
		if keepers[label] {
			pruned = append(pruned, label)
		}
	}
	return pruned
}

func openExcelInput(filePath string) (*excelize.File, error) {
	inputFile, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error opening input file %s: %w", filePath, err)
	}
	return inputFile, nil
}
