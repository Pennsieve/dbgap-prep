package sampleattributes

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHeaderToAttributeLabels(t *testing.T) {
	var samplesHeader = []string{"sample id", "subject id", "was derived from", "pool id", "sample experimental group", "sample type", "sample anatomical location", "also in dataset", "member of", "metadata only", "number of directly derived samples", "laboratory internal id", "date of derivation", "experimental log file path", "reference atlas", "pathology", "laterality", "cell type", "plane of section", "protocol title", "protocol url or doi", "RNA concentration", "RNA concentration method", "RNA purity", "RNA quality", "RNA quality method", "amputation", "body temp", "collection institution", "cross clamp (first)", "cross clamp (last)", "cross clamp time", "donor status", "fixation method", "fixation temp", "fixation time", "freeze thaw cycles", "freezing method", "freezing temp", "incision time", "ischemic time", "post-mortem interval", "protein concentration", "protein concentration method", "sample collection site", "storage temp", "time of sample collection", "nCells", "nFeature_RNA", "nCount_RNA", "nFeature_ATAC", "nCount_ATAC"}

	attrLabels := HeaderToAttributeLabels(samplesHeader)

	assert.Len(t, attrLabels, len(samplesHeader)-2)
	assert.Equal(t, samplesHeader[2:], attrLabels)
}
