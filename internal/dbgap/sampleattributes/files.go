package sampleattributes

import (
	"fmt"
	"strings"

	"github.com/pennsieve/dbgap-prep/internal/dbgap/dd"
	"github.com/pennsieve/dbgap-prep/internal/dbgap/ds"
	sampleattributesdd "github.com/pennsieve/dbgap-prep/internal/dbgap/sampleattributes/dd"
	sampleattributesds "github.com/pennsieve/dbgap-prep/internal/dbgap/sampleattributes/ds"
	"github.com/pennsieve/dbgap-prep/internal/dbgap/sampleattributes/models"
	"github.com/pennsieve/dbgap-prep/internal/enums/analytetype"
	"github.com/pennsieve/dbgap-prep/internal/enums/istumor"

	"github.com/pennsieve/dbgap-prep/internal/samples"
)

var canonicalVariables = []dd.Variable{*dd.SampleIDVar, models.BodySiteVar, models.AnalyteTypeVar, models.IsTumorVar, models.LateralityVar, models.SampleCollectionSiteVar, models.SPARCDatasetDOIVar}

func WriteFiles(outputDirectory string, analyteType analytetype.Type, isTumor istumor.Value, samplesHeader []string, consentedSubjectSamples []samples.Sample) error {
	variables := presentVariables(samplesHeader)
	spec := sampleattributesdd.Spec(variables)
	ddWriter := dd.NewXLSXWriter(outputDirectory, spec.FileName)

	if err := ddWriter.Write(spec); err != nil {
		return fmt.Errorf("error writing sample attributes file: %w", err)
	}

	dsWriter := ds.NewXLSXWriter(outputDirectory, sampleattributesds.DefaultFileNameBase)
	if err := sampleattributesds.Write(dsWriter, analyteType, isTumor, variables, consentedSubjectSamples); err != nil {
		return fmt.Errorf("error writing sample attributes file: %w", err)
	}

	return nil
}

func presentVariables(samplesHeader []string) []dd.Variable {
	var headerSet = make(map[string]bool, len(samplesHeader))
	for _, header := range samplesHeader {
		headerSet[strings.ToLower(header)] = true
	}
	present := make([]dd.Variable, 0, len(canonicalVariables))
	for _, canonicalVariable := range canonicalVariables {
		if headerSet[strings.ToLower(canonicalVariable.SourceColumnName)] || canonicalVariable.Name == models.IsTumorVar.Name || canonicalVariable.Name == models.AnalyteTypeVar.Name || canonicalVariable.Name == models.SPARCDatasetDOIVar.Name {
			present = append(present, canonicalVariable)
		}
	}
	return present
}
