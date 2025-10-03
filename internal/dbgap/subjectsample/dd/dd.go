package dd

import (
	"github.com/pennsieve/dbgap-prep/internal/dbgap/dd"
)

var Spec = dd.Spec{
	FileName:  "3b_SSM_DD.xlsx",
	SheetName: "3b_SSM_DD",
	Rows: [][]any{
		{dd.VarNameColumn, dd.VarDescColumn, dd.TypeColumn, dd.ValuesColumn},
		dd.SubjectIDVar.ToDDRow(),
		dd.SampleIDVar.ToDDRow(),
	},
}
