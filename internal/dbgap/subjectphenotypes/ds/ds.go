package ds

import (
	"github.com/pennsieve/dbgap-prep/internal/dbgap/dd"
	"github.com/pennsieve/dbgap-prep/internal/dbgap/ds"
	"github.com/pennsieve/dbgap-prep/internal/subjects"
)

const FileName = "5a_SubjectPhenotypes_DS.txt"

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

func Write(path string, variables []dd.Variable, consentedSubjects []subjects.Subject) error {
	variableNames := dd.VariableNames(variables)
	rows := ds.ToRows(variableNames, consentedSubjects, ToRow)

	spec := ds.Spec{FileName: FileName, Header: variableNames}
	return ds.Write(path, spec, rows)
}
