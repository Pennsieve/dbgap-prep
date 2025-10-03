package subjectconsent

import (
	"github.com/pennsieve/dbgap-prep/internal/dbgap"
	"github.com/pennsieve/dbgap-prep/internal/dbgap/dd"
	"github.com/pennsieve/dbgap-prep/internal/subjects"
)

var GRUConsent = dd.NewEncodedValue("1", "General Research Use (GRU)")

var ConsentVar = dbgap.Variable{
	Name:        "CONSENT",
	Description: "Registered consent groups (Data Use Limitations (DUL)) as determined by submitters' Institutional Review Boards (IRB) or equivalent body.",
	Type:        dbgap.EncodedValueType,
	Values:      []dd.EncodedValue{GRUConsent},
}

// ConsentFromSubject always returns GRUConsent until we learn otherwise.
func ConsentFromSubject(_ subjects.Subject) string {
	return GRUConsent.Value
}
