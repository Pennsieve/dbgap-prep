package models

import (
	"github.com/pennsieve/dbgap-prep/internal/dbgap/dd"
	"github.com/pennsieve/dbgap-prep/internal/subjects"
)

// NoConsent should not be one of the Values in the DD file per dbGaP documentation.
// We just use it internally to check if a subject is consented or not. But these things
// are in flux.
var NoConsent = dd.NewEncodedValue("0", "No Consent (NC)")

var GRUConsent = dd.NewEncodedValue("1", "General Research Use (GRU)")

var ConsentVar = &dd.Variable{
	Name:        "CONSENT",
	Description: "Registered consent groups (Data Use Limitations (DUL)) as determined by submitters' Institutional Review Boards (IRB) or equivalent body.",
	Type:        dd.EncodedValueType,
	Values:      []dd.EncodedValue{GRUConsent},
}

// ConsentFromSubject always returns GRUConsent until we learn otherwise.
func ConsentFromSubject(_ subjects.Subject) string {
	return GRUConsent.Value
}
