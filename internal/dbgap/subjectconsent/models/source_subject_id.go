package models

import (
	"github.com/pennsieve/dbgap-prep/internal/dbgap/dd"
	"github.com/pennsieve/dbgap-prep/internal/subjects"
)

var SourceSubjectID = &dd.Variable{
	Name:             "SOURCE_SUBJECT_ID",
	Description:      "Subject identifier used by the external source repository named in SUBJECT_SOURCE, linking a SPARC subject to its corresponding dbGaP submission record, formatted as <SPARC study UUID>#sub-XXXX",
	Type:             dd.StringType,
	SourceColumnName: subjects.SourceSubjectIDdbGapColumn,
}

func SourceSubjectIDFromSubject(subject subjects.Subject) string {
	return subject.SourceSubjectIDdbGap
}
