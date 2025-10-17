package models

import "github.com/pennsieve/dbgap-prep/internal/dbgap/dd"

var BodySiteVar = &dd.Variable{
	Name:        "BODY_SITE",
	Description: "Body site where sample was collected",
	Type:        dd.StringType,
}

var AnalyteTypeVar = &dd.Variable{
	Name:        "ANALYTE_TYPE",
	Description: "Analyte type",
	Type:        dd.StringType,
}

var IsTumor = dd.NewEncodedValue("Y", "Is tumor")
var NotTumor = dd.NewEncodedValue("N", "Is not a tumor")
var UnknownTumor = dd.NewEncodedValue("UNK", "Tumor status unknown")

var IsTumorVar = &dd.Variable{
	Name:        "IS_TUMOR",
	Description: "Tumor status",
	Type:        dd.EncodedValueType,
	Values:      []dd.EncodedValue{IsTumor, NotTumor, UnknownTumor},
}
