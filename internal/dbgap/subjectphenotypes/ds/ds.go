package ds

import (
	"github.com/pennsieve/dbgap-prep/internal/dbgap/dd"
	"github.com/pennsieve/dbgap-prep/internal/dbgap/ds"
	"github.com/pennsieve/dbgap-prep/internal/subjects"
)

const DefaultFileNameBase = "5a_SubjectPhenotypes_DS"
const FileName = "5a_SubjectPhenotypes_DS.xlsx"

func ToRow(variableNames []string, subject subjects.Subject) []string {
	row := make([]string, 0, len(variableNames))
	for _, variableName := range variableNames {
		var value string
		switch variableName {
		case dd.SubjectIDVar.Name:
			value = subject.ID
		default:
			value = subject.Values[variableName]
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
