package ds

import (
	"github.com/pennsieve/dbgap-prep/internal/dbgap/dd"
	"github.com/pennsieve/dbgap-prep/internal/dbgap/ds"
	"github.com/pennsieve/dbgap-prep/internal/samples"
)

const DefaultFileNameBase = "6a_SampleAttributes_DS"

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

func Write(writer ds.Writer, variables []dd.Variable, consentedSubjectSamples []samples.Sample) error {
	rows := ds.ToRows(variables, consentedSubjectSamples, ToRow)

	spec := ds.Spec{Variables: variables}
	return writer.Write(spec, rows)
}
