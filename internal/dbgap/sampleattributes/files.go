package sampleattributes

import (
	"github.com/pennsieve/dbgap-prep/internal/dbgap/dd"
	"github.com/pennsieve/dbgap-prep/internal/dbgap/ds"
	sampleattributesdd "github.com/pennsieve/dbgap-prep/internal/dbgap/sampleattributes/dd"
	sampleattributesds "github.com/pennsieve/dbgap-prep/internal/dbgap/sampleattributes/ds"

	"github.com/pennsieve/dbgap-prep/internal/logging"
	"github.com/pennsieve/dbgap-prep/internal/samples"
	"log/slog"
	"path/filepath"
)

var logger = logging.PackageLogger("sampleattributes")

func WriteFiles(outputDirectory string, samplesHeader []string, consentedSubjectSamples []samples.Sample) error {
	attrLabels := HeaderToAttributeLabels(samplesHeader)
	variables := sampleattributesdd.Variables(attrLabels)
	spec := sampleattributesdd.Spec(variables)
	sampleAttributesDDPath := filepath.Join(outputDirectory, spec.FileName)

	if err := dd.Write(sampleAttributesDDPath, spec); err != nil {
		return err
	}

	logger.Info("wrote sample attributes DD file", slog.String("file", sampleAttributesDDPath))

	dsWriter := ds.NewXLSXWriter(outputDirectory, sampleattributesds.DefaultFileNameBase)
	if err := sampleattributesds.Write(dsWriter, variables, consentedSubjectSamples); err != nil {
		return err
	}

	logger.Info("wrote sample attributes DS file", slog.String("file", dsWriter.Path()))

	return nil
}

func HeaderToAttributeLabels(samplesHeader []string) []string {
	cleaned := make([]string, 0, len(samplesHeader))
	for _, label := range samplesHeader {
		if label != samples.IDLabel && label != samples.SubjectIDLabel {
			cleaned = append(cleaned, label)
		}
	}
	return cleaned
}
