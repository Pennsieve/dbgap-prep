package dd

import (
	"github.com/pennsieve/dbgap-prep/internal/dbgap"
	"github.com/pennsieve/dbgap-prep/internal/dbgap/dd"
	"github.com/pennsieve/dbgap-prep/internal/dbgap/subjectconsent"
)

var Spec = dd.Spec{
	FileName:  "2b_SubjectConsent_DD.xlsx",
	SheetName: "2b_SubjectConsent_DD",
	Rows: [][]any{
		{dd.VarNameColumn, dd.VarDescColumn, dd.TypeColumn, dd.ValuesColumn},
		dbgap.SubjectIDVar.ToDDRow(),
		subjectconsent.ConsentVar.ToDDRow(),
		subjectconsent.SexVar.ToDDRow(),
	},
}
