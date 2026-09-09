package models

import (
	"github.com/pennsieve/dbgap-prep/internal/dbgap/dd"
	"github.com/pennsieve/dbgap-prep/internal/subjects"
)

const SPARCSubjectSource = "SPARC.SCIENCE"

var SubjectSourceVar = &dd.Variable{
	Name:        "SUBJECT_SOURCE",
	Description: "Name of the external repository or source that holds complementary data for the subject.",
	Type:        dd.StringType,
}

func SubjectSourceFromSubject(_ subjects.Subject) string {
	return SPARCSubjectSource
}
