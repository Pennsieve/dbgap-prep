package ds

import (
	"github.com/pennsieve/dbgap-prep/internal/config"
	"github.com/pennsieve/dbgap-prep/internal/dbgap/dd"
	"github.com/pennsieve/dbgap-prep/internal/dbgap/ds"
	"github.com/pennsieve/dbgap-prep/internal/dbgap/sampleattributes/models"
	"github.com/pennsieve/dbgap-prep/internal/samples"
)

const DefaultFileNameBase = "6a_SampleAttributes_DS"

func NewToRow(analyteType config.AnalyteType, isTumor bool) ds.ToRowFunc[samples.Sample] {
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
				value = models.ToIsTumorValue(isTumor)
			default:
				value = sample.Values[variable.SourceColumnName]
			}
			row = append(row, value)
		}
		return row
	}
}

func Write(writer ds.Writer, analyteType config.AnalyteType, isTumor bool, variables []dd.Variable, consentedSubjectSamples []samples.Sample) error {
	rows := ds.ToRows(variables, consentedSubjectSamples, NewToRow(analyteType, isTumor))

	spec := ds.Spec{Variables: variables}
	return writer.Write(spec, rows)
}
