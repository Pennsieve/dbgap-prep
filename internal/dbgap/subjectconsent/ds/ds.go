package ds

import (
	"github.com/pennsieve/dbgap-prep/internal/dbgap/dd"
	"github.com/pennsieve/dbgap-prep/internal/dbgap/ds"
	"github.com/pennsieve/dbgap-prep/internal/dbgap/subjectconsent/models"
	"github.com/pennsieve/dbgap-prep/internal/subjects"
)

var Spec = ds.Spec{
	FileName: "2a_SubjectConsent_DS.txt",
	Header:   []string{dd.SubjectIDVar.Name, models.ConsentVar.Name, models.SexVar.Name},
}

type SubjectConsent struct {
	SubjectID string
	Consent   string
	Sex       string
}

func (sc SubjectConsent) IsConsented() bool {
	return len(sc.Consent) == 0 || models.NoConsent.Value != sc.Consent
}

func ToDSRow(subject subjects.Subject) ([]string, SubjectConsent) {
	// Order of items in slice must match the Header row
	subjectConsent := SubjectConsent{
		SubjectID: subject.ID,
		Consent:   models.ConsentFromSubject(subject),
		Sex:       models.SexFromSubject(subject),
	}
	row := []string{
		subjectConsent.SubjectID,
		subjectConsent.Consent,
		subjectConsent.Sex,
	}
	return row, subjectConsent
}

func ToDSRows(subs []subjects.Subject) ([][]string, []SubjectConsent) {
	rows := make([][]string, 0, len(subs))
	subjectConsents := make([]SubjectConsent, 0, len(subs))
	for _, subject := range subs {
		row, sc := ToDSRow(subject)
		rows = append(rows, row)
		subjectConsents = append(subjectConsents, sc)
	}
	return rows, subjectConsents
}

func Write(path string, subs []subjects.Subject) ([]SubjectConsent, error) {
	rows, subjectConsents := ToDSRows(subs)
	err := ds.Write(path, Spec, rows)
	if err != nil {
		return nil, err
	}
	return subjectConsents, nil
}
