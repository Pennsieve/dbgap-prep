package models

import (
	"fmt"

	"github.com/pennsieve/dbgap-prep/internal/dbgap/dd"
	"github.com/pennsieve/dbgap-prep/internal/enums/consentgroup"
	"github.com/pennsieve/dbgap-prep/internal/subjects"
)

// NoConsent should not be one of the Values in the DD file per dbGaP documentation.
// We just use it internally to check if a subject is consented or not. But these things
// are in flux.
var NoConsent = dd.NewEncodedValue("0", "No Consent (NC)")

var ConsentedValue = "1"

var GRUConsent = dd.NewEncodedValue(ConsentedValue, "General Research Use (GRU)")

var HMBConsent = dd.NewEncodedValue(ConsentedValue, "Health/Medical/Biomedical (HMB)")

var OtherConsent = dd.NewEncodedValue(ConsentedValue, "Other (PLEASE CHANGE)")

func ConsentVariable(consentGroup consentgroup.Group) (dd.Variable, error) {
	consentVariable := dd.Variable{
		Name:        "CONSENT",
		Description: "Consent group as determined by DAC",
		Type:        dd.EncodedValueType,
	}
	switch consentGroup {
	case consentgroup.GRU:
		consentVariable.Values = []dd.EncodedValue{GRUConsent}
	case consentgroup.HMB:
		consentVariable.Values = []dd.EncodedValue{HMBConsent}
	case consentgroup.OTHER:
		consentVariable.Values = []dd.EncodedValue{OtherConsent}
	default:
		return dd.Variable{}, fmt.Errorf("unknown consent group: %d", consentGroup)

	}
	return consentVariable, nil
}

// ConsentFromSubject always returns ConsentedValue
func ConsentFromSubject(_ subjects.Subject) string {
	return ConsentedValue
}
