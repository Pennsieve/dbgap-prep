package ds

import (
	"github.com/pennsieve/dbgap-prep/internal/dbgap/dd"
	"github.com/pennsieve/dbgap-prep/internal/dbgap/ds"
	"github.com/pennsieve/dbgap-prep/internal/dbgap/sampleattributes/models"
	"github.com/pennsieve/dbgap-prep/internal/enums/analytetype"
	"github.com/pennsieve/dbgap-prep/internal/enums/istumor"
	"github.com/pennsieve/dbgap-prep/internal/samples"
)

const DefaultFileNameBase = "6a_SampleAttributes_DS"

func NewToRow(analyteType analytetype.Type, isTumor istumor.Value) ds.ToRowFunc[samples.Sample] {
	return func(variables []dd.Variable, sample samples.Sample) []string {
		row := make([]string, 0, len(variables))
		for _, variable := range variables {
			var value string
			switch variable.Name {
			case dd.SampleIDVar.Name:
				value = sample.ID
			case models.AnalyteTypeVar.Name:
				value = analyteType.String()
			case models.IsTumorVar.Name:
				value = isTumor.String()
			default:
				value = sample.Values[variable.SourceColumnName]
			}
			row = append(row, value)
		}
		return row
	}
}

func Write(writer ds.Writer, analyteType analytetype.Type, isTumor istumor.Value, variables []dd.Variable, consentedSubjectSamples []samples.Sample) error {
	rows := ds.ToRows(variables, consentedSubjectSamples, NewToRow(analyteType, isTumor))

	spec := ds.Spec{Variables: variables}
	return writer.Write(spec, rows)
}
