package subjectphenotypes

import (
	"fmt"
	"strings"

	"github.com/pennsieve/dbgap-prep/internal/dbgap/dd"
	"github.com/pennsieve/dbgap-prep/internal/dbgap/ds"
	subjectgphenotypesdd "github.com/pennsieve/dbgap-prep/internal/dbgap/subjectphenotypes/dd"
	subjectgphenotypesds "github.com/pennsieve/dbgap-prep/internal/dbgap/subjectphenotypes/ds"
	"github.com/pennsieve/dbgap-prep/internal/dbgap/subjectphenotypes/models"
	"github.com/pennsieve/dbgap-prep/internal/subjects"
)

var canonicalVariables = []dd.Variable{*dd.SubjectIDVar, models.AgeVariable, models.SexVariable, models.RaceVariable, models.SubjectExperimentalGroupVariable}

func WriteFiles(outputDirectory string, subjectsHeader []string, consentedSubjects []subjects.Subject) error {
	variables := presentVariables(subjectsHeader)
	spec := subjectgphenotypesdd.Spec(variables)

	ddWriter := dd.NewXLSXWriter(outputDirectory, spec.FileName)

	if err := ddWriter.Write(spec); err != nil {
		return fmt.Errorf("error writing subject phenotypes file: %w", err)
	}

	dsWriter := ds.NewXLSXWriter(outputDirectory, subjectgphenotypesds.DefaultFileNameBase)
	if err := subjectgphenotypesds.Write(dsWriter, variables, consentedSubjects); err != nil {
		return fmt.Errorf("error writing subject phenotypes file: %w", err)
	}

	return nil
}

func presentVariables(subjectsHeader []string) []dd.Variable {
	var headerSet = make(map[string]bool, len(subjectsHeader))
	for _, header := range subjectsHeader {
		headerSet[strings.ToLower(header)] = true
	}
	present := make([]dd.Variable, 0, len(canonicalVariables))
	for _, canonicalVariable := range canonicalVariables {
		if headerSet[strings.ToLower(canonicalVariable.SourceColumnName)] {
			present = append(present, canonicalVariable)
		}
	}
	return present
}
