package dd

import (
	"github.com/pennsieve/dbgap-prep/internal/dbgap/dd"
	"github.com/pennsieve/dbgap-prep/internal/dbgap/subjectconsent/models"
)

var header = []dd.Column{dd.VarNameColumn, dd.VarDescColumn, dd.TypeColumn, dd.ValuesColumn}

func Spec(consentVariable *dd.Variable) dd.Spec {
	spec := dd.Spec{
		FileName:  "2b_SubjectConsent_DD.xlsx",
		SheetName: "2b_SubjectConsent_DD",
		Header:    header,
		Rows: [][]any{
			dd.SubjectIDVar.ToDDRow(header),
			consentVariable.ToDDRow(header),
			models.SexVar.ToDDRow(header),
			models.SubjectSourceVar.ToDDRow(header),
			models.SourceSubjectID.ToDDRow(header),
		},
	}
	return spec
}
