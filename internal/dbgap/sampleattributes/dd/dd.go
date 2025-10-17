package dd

import (
	"github.com/pennsieve/dbgap-prep/internal/dbgap/dd"
	"github.com/pennsieve/dbgap-prep/internal/dbgap/sampleattributes/models"
)

var header = []dd.Column{dd.VarNameColumn, dd.VarDescColumn, dd.TypeColumn, dd.UniqueKeyColumn, dd.ValuesColumn}

func Spec(variables []dd.Variable) dd.Spec {
	rows := make([][]any, 0, len(variables))
	for _, variable := range variables {
		rows = append(rows, variable.ToDDRow(header))
	}
	return dd.Spec{
		FileName:  "6b_SampleAttributes_DD.xlsx",
		SheetName: "6b_SampleAttributes_DD",
		Header:    header,
		Rows:      rows,
	}
}

func Variables(attrLabels []string) []dd.Variable {
	variables := make([]dd.Variable, 4, len(attrLabels)+4)
	variables[0] = *dd.SampleIDVar.With(dd.UniqueKeyColumn, "X")
	variables[1] = *models.BodySiteVar
	variables[2] = *models.AnalyteTypeVar
	variables[3] = *models.IsTumorVar

	//TODO figure out what to actually do here.
	for _, label := range attrLabels {
		variables = append(variables, dd.Variable{
			Name:        label,
			Description: label,
			Type:        dd.StringType,
		})
	}
	return variables
}
