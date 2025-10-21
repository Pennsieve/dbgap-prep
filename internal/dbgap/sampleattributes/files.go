package sampleattributes

import (
	"fmt"
	"github.com/pennsieve/dbgap-prep/internal/dbgap/dd"
	"github.com/pennsieve/dbgap-prep/internal/dbgap/ds"
	sampleattributesdd "github.com/pennsieve/dbgap-prep/internal/dbgap/sampleattributes/dd"
	sampleattributesds "github.com/pennsieve/dbgap-prep/internal/dbgap/sampleattributes/ds"

	"github.com/pennsieve/dbgap-prep/internal/samples"
)

func WriteFiles(outputDirectory string, samplesHeader []string, consentedSubjectSamples []samples.Sample) error {
	attrLabels := HeaderToAttributeLabels(samplesHeader)
	variables := sampleattributesdd.Variables(attrLabels)
	spec := sampleattributesdd.Spec(variables)
	ddWriter := dd.NewNoOpWriter(outputDirectory, spec.FileName)

	if err := ddWriter.Write(spec); err != nil {
		return fmt.Errorf("error writing sample attributes file: %w", err)
	}

	dsWriter := ds.NewXLSXWriter(outputDirectory, sampleattributesds.DefaultFileNameBase)
	if err := sampleattributesds.Write(dsWriter, variables, consentedSubjectSamples); err != nil {
		return fmt.Errorf("error writing sample attributes file: %w", err)
	}

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
