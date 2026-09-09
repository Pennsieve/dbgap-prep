package models

import (
	"github.com/pennsieve/dbgap-prep/internal/dbgap/dd"
	"github.com/pennsieve/dbgap-prep/internal/subjects"
)

// SubjectExperimentalGroupVariableName is the name of the variable in the dbGaP files we are generating
const SubjectExperimentalGroupVariableName = "SUBJECT_EXPERIMENTAL_GROUP"

var SubjectExperimentalGroupVariable = dd.Variable{
	Name:             SubjectExperimentalGroupVariableName,
	Description:      "Experimental group assignment of the subject",
	Type:             dd.StringType,
	SourceColumnName: "subject experimental group",
}

func SubjectExperimentalGroupFromSubject(subject subjects.Subject) string {
	subjectExperimentalGroup, _ := subject.SearchValue(SubjectExperimentalGroupVariable.SourceColumnName)
	return subjectExperimentalGroup
}
