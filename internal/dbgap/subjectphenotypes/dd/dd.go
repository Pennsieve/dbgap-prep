package dd

import (
	"github.com/pennsieve/dbgap-prep/internal/dbgap/dd"
)

var header = []dd.Column{dd.VarNameColumn, dd.VarDescColumn, dd.TypeColumn, dd.UniqueKeyColumn, dd.ValuesColumn}

func Spec(variables []dd.Variable) dd.Spec {
	rows := make([][]any, 0, len(variables))
	for _, variable := range variables {
		rows = append(rows, variable.ToDDRow(header))
	}
	return dd.Spec{
		FileName:  "5b_SubjectPhenotypes_DD.xlsx",
		SheetName: "5b_SubjectPhenotypes_DD",
		Header:    header,
		Rows:      rows,
	}
}

func Variables(attrLabels []string) []dd.Variable {
	variables := make([]dd.Variable, 1, len(attrLabels)+1)
	variables[0] = *dd.SubjectIDVar.With(dd.UniqueKeyColumn, "X")

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
