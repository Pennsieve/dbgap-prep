package ds

import (
	"github.com/pennsieve/dbgap-prep/internal/dbgap/dd"
	"github.com/pennsieve/dbgap-prep/internal/dbgap/ds"
	"github.com/pennsieve/dbgap-prep/internal/dbgap/subjectconsent/models"
	"github.com/pennsieve/dbgap-prep/internal/samples"
	"github.com/pennsieve/dbgap-prep/internal/subjects"
)

var Spec = ds.Spec{
	FileName:  "2a_SubjectConsent_DS.txt",
	Variables: []dd.Variable{*dd.SubjectIDVar, *models.ConsentVar, *models.SexVar},
}

type SubjectConsent struct {
	SubjectID string
	Consent   string
	Sex       string
}

func (sc SubjectConsent) IsConsented() bool {
	return len(sc.Consent) == 0 || models.NoConsent.Value != sc.Consent
}

func ToRow(subject subjects.Subject) ([]string, SubjectConsent) {
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

func ToRows(subs []subjects.Subject) ([][]string, []SubjectConsent) {
	rows := make([][]string, 0, len(subs))
	subjectConsents := make([]SubjectConsent, 0, len(subs))
	for _, subject := range subs {
		row, sc := ToRow(subject)
		rows = append(rows, row)
		subjectConsents = append(subjectConsents, sc)
	}
	return rows, subjectConsents
}

func Write(path string, subs []subjects.Subject) ([]SubjectConsent, error) {
	rows, subjectConsents := ToRows(subs)
	err := ds.Write(path, Spec, rows)
	if err != nil {
		return nil, err
	}
	return subjectConsents, nil
}

// GetConsented returns three slices, the first is those subjects.Subject that are consented (i.e., consent != "0").
// The second slice are those samples.Sample that have a consented subject in their original samples.xlsx order.
// The third slice are those same consented samples, but in the order their subjects appear in subjects.xlsx.
func GetConsented(subjectConsents []SubjectConsent, subs []subjects.Subject, samps []samples.Sample) ([]subjects.Subject, []samples.Sample, []samples.Sample) {
	consentsBySubjectID := make(map[string]bool, len(subjectConsents))
	for _, consent := range subjectConsents {
		consentsBySubjectID[consent.SubjectID] = consent.IsConsented()
	}

	consentedSubjectToPosition := make(map[string]int, len(subs))

	// Only include subjects with consent
	consentedSubjects := make([]subjects.Subject, 0, len(subs))
	for i, subject := range subs {
		if consentsBySubjectID[subject.ID] {
			consentedSubjects = append(consentedSubjects, subject)
			consentedSubjectToPosition[subject.ID] = i
		}
	}

	consentedSamples := make([]samples.Sample, 0, len(samps))
	samplesBySubjectPosition := make([][]samples.Sample, len(subjectConsents))
	for _, sample := range samps {
		if position, isConsented := consentedSubjectToPosition[sample.SubjectID]; isConsented {
			consentedSamples = append(consentedSamples, sample)
			samplesBySubjectPosition[position] = append(samplesBySubjectPosition[position], sample)
		}
	}

	// Flatten the samples grouped by subject position to a slice of samples.
	consentedSamplesOrderedBySubject := make([]samples.Sample, 0, len(consentedSamples))
	for _, subjectSamples := range samplesBySubjectPosition {
		consentedSamplesOrderedBySubject = append(consentedSamplesOrderedBySubject, subjectSamples...)
	}

	return consentedSubjects, consentedSamples, consentedSamplesOrderedBySubject
}
