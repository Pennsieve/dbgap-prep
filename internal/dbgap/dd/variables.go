package dd

import (
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
	Attributes  map[Column]any
	// Values should be left nil (or empty) if Type does not include EncodedValueType
	Values []EncodedValue
}

// ToDDRow returns this variable's attribute values as a slice in the same order
// as the passed in header. If ValuesColumn appears in header, it must be the final element.
func (v *Variable) ToDDRow(header []Column) []any {
	var row []any
	for _, c := range header {
		switch c {
		case VarNameColumn:
			row = append(row, v.Name)
		case VarDescColumn:
			row = append(row, v.Description)
		case TypeColumn:
			row = append(row, v.Type)
		case ValuesColumn:
			for _, value := range v.Values {
				row = append(row, value)
			}
		default:
			if attr, ok := v.Attributes[c]; ok {
				row = append(row, attr)
			} else {
				row = append(row, "")
			}

		}
	}
	return row
}

func (v *Variable) With(column Column, value any) *Variable {
	if v.Attributes == nil {
		v.Attributes = map[Column]any{}
	}
	v.Attributes[column] = value
	return v
}

var SubjectIDVar = &Variable{
	Name:        "SUBJECT_ID",
	Description: "Subject ID",
	Type:        StringType,
}

var SampleIDVar = &Variable{
	Name:        "SAMPLE_ID",
	Description: "Sample ID",
	Type:        StringType,
}

func VariableNames(variables []Variable) []string {
	names := make([]string, 0, len(variables))
	for _, variable := range variables {
		names = append(names, variable.Name)
	}
	return names
}
