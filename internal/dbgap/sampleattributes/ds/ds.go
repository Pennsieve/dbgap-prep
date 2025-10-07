package ds

import (
	"github.com/pennsieve/dbgap-prep/internal/dbgap/dd"
	"github.com/pennsieve/dbgap-prep/internal/dbgap/ds"
	"github.com/pennsieve/dbgap-prep/internal/samples"
)

const FileName = "6a_SampleAttributes_DD.txt"

func ToRow(variableNames []string, sample samples.Sample) []string {
	row := make([]string, 0, len(variableNames))
	for _, variableName := range variableNames {
		var value string
		switch variableName {
		case dd.SampleIDVar.Name:
			value = sample.ID
		default:
			value = sample.Values[variableName]
		}
		row = append(row, value)
	}
	return row
}

func VariableNames(variables []dd.Variable) []string {
	names := make([]string, 0, len(variables))
	for _, variable := range variables {
		names = append(names, variable.Name)
	}
	return names
}

func ToRows(variableNames []string, consentedSubjectSamples map[string][]samples.Sample) [][]string {
	var rows [][]string
	for _, consentedSamples := range consentedSubjectSamples {
		for _, consentedSample := range consentedSamples {
			rows = append(rows, ToRow(variableNames, consentedSample))
		}
	}
	return rows
}

func Write(path string, variables []dd.Variable, consentedSubjectSamples map[string][]samples.Sample) error {
	variableNames := VariableNames(variables)
	spec := ds.Spec{FileName: FileName, Header: variableNames}
	rows := ToRows(variableNames, consentedSubjectSamples)
	return ds.Write(path, spec, rows)
}
