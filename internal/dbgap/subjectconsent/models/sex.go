package models

import (
	"github.com/pennsieve/dbgap-prep/internal/dbgap/dd"
	"github.com/pennsieve/dbgap-prep/internal/subjects"
	"strings"
)

var MaleSex = dd.NewEncodedValue("1", "Male")
var FemaleSex = dd.NewEncodedValue("2", "Female")
var UnknownSex = dd.NewEncodedValue("UNK", "Unknown")

var SexVar = &dd.Variable{
	Name:        "SEX",
	Description: "Biological Sex",
	Type:        dd.EncodedValueType,
	Values:      []dd.EncodedValue{MaleSex, FemaleSex, UnknownSex},
}

func SexFromSubject(subject subjects.Subject) string {
	// subjects file uses values 1,2,3,4 for biological sex, but dbGaP wants 1, 2, or UNK
	switch strings.ToLower(subject.Sex) {
	case "male":
		return MaleSex.Value
	case "female":
		return FemaleSex.Value
	default:
		return UnknownSex.Value

	}
}
