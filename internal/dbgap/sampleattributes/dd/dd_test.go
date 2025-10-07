package dd

import (
	"github.com/pennsieve/dbgap-prep/internal/dbgap/dd"
	"github.com/pennsieve/dbgap-prep/internal/dbgap/sampleattributes/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestVariables(t *testing.T) {
	var attributeLabels = []string{"was derived from", "pool id", "sample experimental group", "sample type", "sample anatomical location", "also in dataset", "member of", "metadata only", "number of directly derived samples", "laboratory internal id", "date of derivation", "experimental log file path", "reference atlas", "pathology", "laterality", "cell type", "plane of section", "protocol title", "protocol url or doi", "RNA concentration", "RNA concentration method", "RNA purity", "RNA quality", "RNA quality method", "amputation", "body temp", "collection institution", "cross clamp (first)", "cross clamp (last)", "cross clamp time", "donor status", "fixation method", "fixation temp", "fixation time", "freeze thaw cycles", "freezing method", "freezing temp", "incision time", "ischemic time", "post-mortem interval", "protein concentration", "protein concentration method", "sample collection site", "storage temp", "time of sample collection", "nCells", "nFeature_RNA", "nCount_RNA", "nFeature_ATAC", "nCount_ATAC"}

	variables := Variables(attributeLabels)

	require.Len(t, variables, len(attributeLabels)+4)

	assert.Equal(t, *dd.SampleIDVar, variables[0])
	assert.Equal(t, "X", variables[0].Attributes[dd.UniqueKeyColumn])

	assert.Equal(t, *models.BodySiteVar, variables[1])

	assert.Equal(t, *models.AnalyteTypeVar, variables[2])

	assert.Equal(t, *models.IsTumorVar, variables[3])

	for i, label := range attributeLabels {
		assert.Equal(t, label, variables[i+4].Name)

	}

}
