package models

import (
	"github.com/pennsieve/dbgap-prep/internal/dbgap/dd"
	"github.com/pennsieve/dbgap-prep/internal/subjects"
)

const SexVariableName = "SEX"

var SexVariable = dd.Variable{
	Name:             SexVariableName,
	Description:      "Biological sex",
	Type:             dd.StringType,
	SourceColumnName: "sex",
}

func SexFromSubject(subject subjects.Subject) string {
	return subject.Sex
}
