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
		FileName:  "6b_SampleAttributes_DD.xlsx",
		SheetName: "6b_SampleAttributes_DD",
		Header:    header,
		Rows:      rows,
	}
}
