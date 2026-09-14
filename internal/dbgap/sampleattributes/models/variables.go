package models

import (
	"github.com/pennsieve/dbgap-prep/internal/dbgap/dd"
)

var BodySiteVar = dd.Variable{
	Name:             "BODY_SITE",
	Description:      "Body site where sample was collected",
	Type:             dd.StringType,
	SourceColumnName: "sample anatomical location",
}

var AnalyteTypeVar = dd.Variable{
	Name:        "ANALYTE_TYPE",
	Description: "Analyte type of the sample",
	Type:        dd.StringType,
}

var IsTumor = dd.NewEncodedValue("Yes", "Is tumor")
var NotTumor = dd.NewEncodedValue("No", "Is not a tumor")

var IsTumorVar = dd.Variable{
	Name:        "IS_TUMOR",
	Description: "Tumor status of the sample",
	Type:        dd.EncodedValueType,
	Values:      []dd.EncodedValue{IsTumor, NotTumor},
}

var LateralityVar = dd.Variable{
	Name:             "LATERALITY",
	Description:      "Side of the body from which the sample was collected",
	Type:             dd.StringType,
	SourceColumnName: "laterality",
}

var SampleCollectionSiteVar = dd.Variable{
	Name:             "SAMPLE_COLLECTION_SITE",
	Description:      "Specific anatomical location where the sample was collected",
	Type:             dd.StringType,
	SourceColumnName: "sample collection site",
}

var SPARCDatasetDOIVar = dd.Variable{
	Name:        "SPARC_DATASET_DOI",
	Description: "DOI of the SPARC dataset holding the open-access counterpart record for this sample",
	Type:        dd.StringType,
}

var CanonicalVariables = []dd.Variable{*dd.SampleIDVar, BodySiteVar, AnalyteTypeVar, IsTumorVar, LateralityVar, SampleCollectionSiteVar, SPARCDatasetDOIVar}
