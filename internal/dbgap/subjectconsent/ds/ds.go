package ds

import (
	"github.com/pennsieve/dbgap-prep/internal/dbgap"
	"github.com/pennsieve/dbgap-prep/internal/dbgap/ds"
	"github.com/pennsieve/dbgap-prep/internal/dbgap/subjectconsent"
	"github.com/pennsieve/dbgap-prep/internal/subjects"
)

var Spec = ds.Spec{
	FileName: "2a_SubjectConsent_DS.txt",
	Header:   []string{dbgap.SubjectIDVar.Name, subjectconsent.ConsentVar.Name, subjectconsent.SexVar.Name},
}

func ToDSRow(subject subjects.Subject) []string {
	// Order of items in slice must match the Header row
	return []string{
		subject.ID,
		subjectconsent.ConsentFromSubject(subject),
		subjectconsent.SexFromSubject(subject)}
}

func ToDSRows(subs []subjects.Subject) [][]string {
	rows := make([][]string, 0, len(subs))
	for _, subject := range subs {
		rows = append(rows, ToDSRow(subject))
	}
	return rows
}

func Write(path string, subs []subjects.Subject) error {
	rows := ToDSRows(subs)
	return ds.Write(path, Spec, rows)
}
