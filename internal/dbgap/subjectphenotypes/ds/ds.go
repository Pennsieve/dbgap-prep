package ds

import (
	"fmt"

	"github.com/pennsieve/dbgap-prep/internal/dbgap/dd"
	"github.com/pennsieve/dbgap-prep/internal/dbgap/ds"
	"github.com/pennsieve/dbgap-prep/internal/dbgap/subjectphenotypes/models"
	"github.com/pennsieve/dbgap-prep/internal/subjects"
)

const DefaultFileNameBase = "5a_SubjectPhenotypes_DS"
const FileName = "5a_SubjectPhenotypes_DS.xlsx"

func ToRow(variables []dd.Variable, subject subjects.Subject) []string {
	row := make([]string, 0, len(variables))
	for _, variable := range variables {
		var value string
		switch variable.Name {
		case dd.SubjectIDVar.Name:
			value = subject.ID
		case models.AgeVariable.Name:
			intAgeMaybe := models.AgeFromSubject(subject)
			if intAgeMaybe != nil {
				value = fmt.Sprintf("%d", *intAgeMaybe)
			}
		case models.SexVariable.Name:
			value = models.SexFromSubject(subject)
		case models.RaceVariable.Name:
			value = models.RaceFromSubject(subject)
		case models.SubjectExperimentalGroupVariable.Name:
			value = models.SubjectExperimentalGroupFromSubject(subject)

		default:
			value = subject.Values[variable.SourceColumnName]
		}
		row = append(row, value)
	}
	return row
}

func Write(writer ds.Writer, variables []dd.Variable, consentedSubjects []subjects.Subject) error {
	rows := ds.ToRows(variables, consentedSubjects, ToRow)

	spec := ds.Spec{Variables: variables}
	return writer.Write(spec, rows)
}
