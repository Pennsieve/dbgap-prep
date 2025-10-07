package dd

import (
	"github.com/pennsieve/dbgap-prep/internal/dbgap/dd"
	"github.com/pennsieve/dbgap-prep/internal/dbgap/subjectconsent/models"
)

var header = []dd.Column{dd.VarNameColumn, dd.VarDescColumn, dd.TypeColumn, dd.ValuesColumn}

var Spec = dd.Spec{
	FileName:  "2b_SubjectConsent_DD.xlsx",
	SheetName: "2b_SubjectConsent_DD",
	Header:    header,
	Rows: [][]any{
		dd.SubjectIDVar.ToDDRow(header),
		models.ConsentVar.ToDDRow(header),
		models.SexVar.ToDDRow(header),
	},
}
