package models

import (
	"github.com/pennsieve/dbgap-prep/internal/dbgap/dd"
	"github.com/pennsieve/dbgap-prep/internal/subjects"
)

const RaceVariableName = "RACE"

var RaceVariable = dd.Variable{
	Name:             RaceVariableName,
	Description:      "Self-reported race",
	Type:             dd.StringType,
	SourceColumnName: "race",
}

func RaceFromSubject(subject subjects.Subject) string {
	race, _ := subject.SearchValue(RaceVariableName)
	return race
}
