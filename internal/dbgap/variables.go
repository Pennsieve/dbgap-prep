package dbgap

import (
	"github.com/pennsieve/dbgap-prep/internal/dbgap/dd"
	"strings"
)

type Type string

const IntegerType = "integer"
const EncodedValueType = "encoded value"
const DecimalType = "decimal"
const StringType = "string"

func MixedType(types []Type) Type {
	strs := make([]string, 0, len(types))
	for _, t := range types {
		strs = append(strs, string(t))
	}
	return Type(strings.Join(strs, ","))
}

type Variable struct {
	Name        string
	Description string
	Type        Type
	// Values should be left nil (or empty) if Type does not include EncodedValueType
	Values []dd.EncodedValue
}

func (v Variable) ToDDRow() []any {
	row := []any{v.Name, v.Description, v.Type}
	for _, value := range v.Values {
		row = append(row, value)
	}
	return row
}

var SubjectIDVar = Variable{
	Name:        "SUBJECT_ID",
	Description: "Subject ID",
	Type:        StringType,
}
