package dd

import (
	"github.com/pennsieve/dbgap-prep/internal/dbgap/dd"
)

var header = []dd.Column{dd.VarNameColumn, dd.VarDescColumn, dd.TypeColumn, dd.ValuesColumn}

var Spec = dd.Spec{
	FileName:  "3b_SSM_DD.xlsx",
	SheetName: "3b_SSM_DD",
	Header:    header,
	Rows: [][]any{
		dd.SubjectIDVar.ToDDRow(header),
		dd.SampleIDVar.ToDDRow(header),
	},
}
