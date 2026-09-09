package models

import (
	"github.com/pennsieve/dbgap-prep/internal/dbgap/dd"
	"github.com/pennsieve/dbgap-prep/internal/subjects"
)

var SourceSubjectID = &dd.Variable{
	Name:        "SOURCE_SUBJECT_ID",
	Description: "Subject identifier used by the external source repository",
	Type:        dd.StringType,
}

func SourceSubjectIDFromSubject(subject subjects.Subject) string {
	return subject.SourceSubjectIDdbGap
}
