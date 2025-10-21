package subjectphenotypes

import (
	"github.com/pennsieve/dbgap-prep/internal/dbgap/dd"
	"github.com/pennsieve/dbgap-prep/internal/dbgap/ds"
	subjectgphenotypesdd "github.com/pennsieve/dbgap-prep/internal/dbgap/subjectphenotypes/dd"
	subjectgphenotypesds "github.com/pennsieve/dbgap-prep/internal/dbgap/subjectphenotypes/ds"
	"github.com/pennsieve/dbgap-prep/internal/subjects"
)

func WriteFiles(outputDirectory string, subjectsHeader []string, consentedSubjects []subjects.Subject) error {
	attrLabels := HeaderToAttributeLabels(subjectsHeader)
	variables := subjectgphenotypesdd.Variables(attrLabels)
	spec := subjectgphenotypesdd.Spec(variables)

	ddWriter := dd.NewNoOpWriter(outputDirectory, spec.FileName)

	if err := ddWriter.Write(spec); err != nil {
		return err
	}

	dsWriter := ds.NewXLSXWriter(outputDirectory, subjectgphenotypesds.DefaultFileNameBase)
	if err := subjectgphenotypesds.Write(dsWriter, variables, consentedSubjects); err != nil {
		return err
	}

	return nil
}

func HeaderToAttributeLabels(subjectsHeader []string) []string {
	cleaned := make([]string, 0, len(subjectsHeader))
	for _, label := range subjectsHeader {
		if label != subjects.IDLabel && label != subjects.SexLabel {
			cleaned = append(cleaned, label)
		}
	}
	return cleaned
}
